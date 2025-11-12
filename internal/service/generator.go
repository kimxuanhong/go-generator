package service

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/xhkzeroone/go-generator/internal/models"
)

type GeneratorService struct {
	manifest *models.Manifest
}

func NewGeneratorService(manifestPath string) (*GeneratorService, error) {
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}
	return &GeneratorService{manifest: manifest}, nil
}

func (s *GeneratorService) GenerateProject(req *GenerateRequest) ([]byte, error) {
	// Validate framework exists
	if _, ok := s.manifest.Frameworks[req.Framework]; !ok {
		return nil, fmt.Errorf("framework '%s' not found", req.Framework)
	}

	// Validate all libs exist
	for _, lib := range req.Libs {
		if _, ok := s.manifest.Libs[lib]; !ok {
			return nil, fmt.Errorf("library '%s' not found", lib)
		}
	}

	// Create temp directory
	tmp, err := os.MkdirTemp("", "gen-"+req.ProjectName+"-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmp)

	// Create base structure
	if err := s.createBaseStructure(tmp, req.IncludeExample, req.Framework); err != nil {
		return nil, fmt.Errorf("failed to create base structure: %w", err)
	}

	// Collect dependencies
	allImports := s.collectDependencies(req)

	// Get framework definition
	fdef := s.manifest.Frameworks[req.Framework]

	// Render framework templates (config.tmpl -> infrastructure, server.tmpl -> not used, handled by app/server.tmpl)
	if err := s.renderFrameworkTemplates(tmp, req, fdef); err != nil {
		return nil, fmt.Errorf("failed to render framework templates: %w", err)
	}

	// Render library templates and merge configs
	includes := make(map[string]bool)
	mergedConfig := make(map[string]interface{})

	// Merge framework config section first
	if fdef.ConfigSection != "" {
		if err := s.mergeConfigSection(fdef.ConfigSection, mergedConfig); err != nil {
			return nil, fmt.Errorf("failed to merge config for framework %s: %w", req.Framework, err)
		}
	}

	// Merge library config sections
	for _, lib := range req.Libs {
		ldef := s.manifest.Libs[lib]
		includes[lib] = true

		if err := s.renderLibTemplates(tmp, req, lib, ldef); err != nil {
			return nil, fmt.Errorf("failed to render library templates for %s: %w", lib, err)
		}

		// Merge config section
		if ldef.ConfigSection != "" {
			if err := s.mergeConfigSection(ldef.ConfigSection, mergedConfig); err != nil {
				return nil, fmt.Errorf("failed to merge config for %s: %w", lib, err)
			}
		}
	}

	// Write config file
	if err := s.writeConfigFile(tmp, mergedConfig); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	// Render Clean Architecture layers only if IncludeExample is true
	if req.IncludeExample {
		// Domain layer
		if err := s.renderDomainLayer(tmp, req, includes); err != nil {
			return nil, fmt.Errorf("failed to render domain layer: %w", err)
		}

		// Repository layer
		if err := s.renderRepositoryLayer(tmp, req, includes); err != nil {
			return nil, fmt.Errorf("failed to render repository layer: %w", err)
		}

		// Usecase layer
		if err := s.renderUsecaseLayer(tmp, req, includes); err != nil {
			return nil, fmt.Errorf("failed to render usecase layer: %w", err)
		}

		// Handler layer
		if err := s.renderHandlerLayer(tmp, req, includes); err != nil {
			return nil, fmt.Errorf("failed to render handler layer: %w", err)
		}

		// Jobs layer (only if cron is included)
		if includes["cron"] {
			if err := s.renderJobsLayer(tmp, req, includes); err != nil {
				return nil, fmt.Errorf("failed to render jobs layer: %w", err)
			}
		}
	}

	// App server (always render, but with or without example routes)
	if err := s.renderAppServer(tmp, req, includes); err != nil {
		return nil, fmt.Errorf("failed to render app server: %w", err)
	}

	// Write main.go
	if err := s.renderMainFile(tmp, req, includes); err != nil {
		return nil, fmt.Errorf("failed to render main file: %w", err)
	}

	// Write deps package
	if err := s.renderDepsPackage(tmp, req, includes); err != nil {
		return nil, fmt.Errorf("failed to render deps package: %w", err)
	}

	// Write go.mod with all dependencies
	if err := s.renderGoMod(tmp, req, allImports); err != nil {
		return nil, fmt.Errorf("failed to render go.mod: %w", err)
	}

	// Render Dockerfile template into project root so generated projects include a Dockerfile
	dockerData := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
		// Binary name defaults to project name
		"BinaryName": req.ProjectName,
		// Default port used by templates; you can change when running the container
		"Port": 8080,
		// Default Go version for builder image
		"GoVersion": "1.20",
	}
	if err := s.renderTemplate("templates/Dockerfile.tmpl", filepath.Join(tmp, "Dockerfile"), dockerData); err != nil {
		return nil, fmt.Errorf("failed to render Dockerfile template: %w", err)
	}

	// Render .gitignore into project root so generated projects include common ignores
	gitignoreData := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
	}
	if err := s.renderTemplate("templates/gitignore.tmpl", filepath.Join(tmp, ".gitignore"), gitignoreData); err != nil {
		return nil, fmt.Errorf("failed to render .gitignore template: %w", err)
	}

	// Render .env.example into project root so generated projects have env var documentation
	envExampleData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	if err := s.renderTemplate("templates/env_example.tmpl", filepath.Join(tmp, ".env.example"), envExampleData); err != nil {
		return nil, fmt.Errorf("failed to render .env.example template: %w", err)
	}

	// Render README.md into project root so generated projects have documentation
	readmeData := map[string]interface{}{
		"ProjectName":    req.ProjectName,
		"ModuleName":     req.ModuleName,
		"Framework":      req.Framework,
		"IncludeExample": req.IncludeExample,
		"Includes":       includes,
	}
	if err := s.renderTemplate("templates/README.tmpl", filepath.Join(tmp, "README.md"), readmeData); err != nil {
		return nil, fmt.Errorf("failed to render README template: %w", err)
	}

	// Create zip file
	zipData, err := s.createZip(tmp)
	if err != nil {
		return nil, fmt.Errorf("failed to create zip: %w", err)
	}

	return zipData, nil
}

func (s *GeneratorService) createBaseStructure(tmp string, includeExample bool, framework string) error {
	dirs := []string{
		"cmd",
		"internal/app",
		"internal/infrastructure",
		"internal/deps",
		"config",
	}

	// Only create example layers if IncludeExample is true
	if includeExample {
		exampleDirs := []string{
			"internal/domain",
			"internal/usecase",
			"internal/repository",
			"internal/adapter/handler",
			"internal/jobs",
		}
		dirs = append(dirs, exampleDirs...)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmp, dir), 0755); err != nil {
			return err
		}
	}
	return nil
}

func (s *GeneratorService) collectDependencies(req *GenerateRequest) []string {
	moduleMap := make(map[string]bool)

	// Add framework imports
	if fdef, ok := s.manifest.Frameworks[req.Framework]; ok {
		for _, imp := range fdef.Imports {
			moduleMap[extractModulePath(imp)] = true
		}
	}

	// Add library imports
	for _, lib := range req.Libs {
		if ldef, ok := s.manifest.Libs[lib]; ok {
			for _, imp := range ldef.Imports {
				moduleMap[extractModulePath(imp)] = true
			}
		}
	}

	// Always add common dependencies
	moduleMap["github.com/sirupsen/logrus"] = true
	moduleMap["github.com/spf13/viper"] = true

	// Convert to sorted slice
	var modules []string
	for mod := range moduleMap {
		modules = append(modules, mod)
	}
	sort.Strings(modules)

	return modules
}

// extractModulePath extracts the module path from an import path
// e.g., "github.com/gofiber/fiber/v2" -> "github.com/gofiber/fiber/v2"
// e.g., "gorm.io/driver/postgres" -> "gorm.io/driver/postgres"
// e.g., "gorm.io/gorm" -> "gorm.io/gorm"
func extractModulePath(importPath string) string {
	// Remove version suffix if present (e.g., /v2, /v3)
	// But keep it if it's part of the module path
	parts := strings.Split(importPath, "/")

	// Check if last part is a version (v1, v2, etc.)
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if strings.HasPrefix(last, "v") && len(last) > 1 {
			// Check if it's a valid version number
			isVersion := true
			for _, r := range last[1:] {
				if r < '0' || r > '9' {
					isVersion = false
					break
				}
			}
			if isVersion {
				// Keep the version in the module path
				return importPath
			}
		}
	}

	// For paths like "gorm.io/driver/postgres", return the first two parts
	// For paths like "github.com/user/repo", return the first three parts
	if strings.Contains(importPath, ".") {
		// Count the number of parts before the first dot
		parts := strings.Split(importPath, "/")
		if len(parts) >= 2 {
			return strings.Join(parts[:min(len(parts), 3)], "/")
		}
	}

	return importPath
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *GeneratorService) renderFrameworkTemplates(tmp string, req *GenerateRequest, fdef models.FrameworkDef) error {
	for _, t := range fdef.Templates {
		var outPath string
		baseName := filepath.Base(strings.TrimSuffix(t, ".tmpl"))

		// render to internal/app/
		outPath = filepath.Join(tmp, "internal/app", baseName+".go")

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		data := map[string]interface{}{
			"ModuleName":  req.ModuleName,
			"ProjectName": req.ProjectName,
			"Framework":   req.Framework,
		}
		if err := s.renderTemplate(t, outPath, data); err != nil {
			return fmt.Errorf("failed to render %s: %w", t, err)
		}
	}
	return nil
}

func (s *GeneratorService) renderLibTemplates(tmp string, req *GenerateRequest, lib string, ldef models.LibDef) error {
	for _, t := range ldef.Templates {
		outPath := filepath.Join(tmp, "internal/infrastructure", lib, filepath.Base(strings.TrimSuffix(t, ".tmpl"))+".go")
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		data := map[string]interface{}{
			"ModuleName":  req.ModuleName,
			"ProjectName": req.ProjectName,
			"Lib":         lib,
		}
		if err := s.renderTemplate(t, outPath, data); err != nil {
			return fmt.Errorf("failed to render %s: %w", t, err)
		}
	}
	return nil
}

func (s *GeneratorService) mergeConfigSection(configPath string, mergedConfig map[string]interface{}) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var section map[string]interface{}
	if err := json.Unmarshal(data, &section); err != nil {
		return err
	}

	for k, v := range section {
		mergedConfig[k] = v
	}

	return nil
}

func (s *GeneratorService) writeConfigFile(tmp string, config map[string]interface{}) error {
	cfgPath := filepath.Join(tmp, "config", "config.json")
	cfgFile, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	enc := json.NewEncoder(cfgFile)
	enc.SetIndent("", "  ")
	return enc.Encode(config)
}

func (s *GeneratorService) renderMainFile(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, "cmd", "main.go")
	data := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
		"Framework":   req.Framework,
		"Includes":    includes,
	}
	return s.renderTemplate("templates/cmd/main.tmpl", outPath, data)
}

func (s *GeneratorService) renderDepsPackage(tmp string, req *GenerateRequest, includes map[string]bool) error {
	// Render deps.go
	depsPath := filepath.Join(tmp, "internal/deps", "deps.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	if err := s.renderTemplate("templates/deps/deps.tmpl", depsPath, data); err != nil {
		return err
	}

	// Render config.go
	configPath := filepath.Join(tmp, "internal/deps", "config.go")
	configData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Framework":  req.Framework,
		"Includes":   includes,
	}
	return s.renderTemplate("templates/deps/config.tmpl", configPath, configData)
}

func (s *GeneratorService) renderGoMod(tmp string, req *GenerateRequest, imports []string) error {
	outPath := filepath.Join(tmp, "go.mod")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Imports":    imports,
	}
	return s.renderTemplate("templates/go_mod.tmpl", outPath, data)
}

func (s *GeneratorService) renderTemplate(tmplPath, dstPath string, data interface{}) error {
	tpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", tmplPath, err)
	}

	return nil
}

func (s *GeneratorService) createZip(tmp string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	err := filepath.WalkDir(tmp, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(tmp, path)
		if err != nil {
			return err
		}

		f, err := zw.Create(rel)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		return err
	})

	if err != nil {
		zw.Close()
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *GeneratorService) renderDomainLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, "internal/domain", "entity.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate("templates/domain/entity.tmpl", outPath, data)
}

func (s *GeneratorService) renderRepositoryLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	// Render user repository
	userRepoPath := filepath.Join(tmp, "internal/repository", "user_repository.go")
	userRepoData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	if err := s.renderTemplate("templates/repository/user_repository.tmpl", userRepoPath, userRepoData); err != nil {
		return err
	}

	// Render cache repository
	cacheRepoPath := filepath.Join(tmp, "internal/repository", "cache_repository.go")
	cacheRepoData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate("templates/repository/cache_repository.tmpl", cacheRepoPath, cacheRepoData)
}

func (s *GeneratorService) renderUsecaseLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, "internal/usecase", "user_usecase.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate("templates/usecase/user_usecase.tmpl", outPath, data)
}

func (s *GeneratorService) renderHandlerLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, "internal/adapter/handler", "user_handler.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Framework":  req.Framework,
		"Includes":   includes,
	}
	return s.renderTemplate("templates/handler/user_handler.tmpl", outPath, data)
}

func (s *GeneratorService) renderAppServer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, "internal/app", "server.go")
	data := map[string]interface{}{
		"ModuleName":     req.ModuleName,
		"ProjectName":    req.ProjectName,
		"Framework":      req.Framework,
		"Includes":       includes,
		"IncludeExample": req.IncludeExample,
	}

	// Use example server template if IncludeExample is true, otherwise use simple server
	templatePath := "templates/app/server_simple.tmpl"
	if req.IncludeExample {
		templatePath = "templates/app/server.tmpl"
	}

	// Render server.go
	if err := s.renderTemplate(templatePath, outPath, data); err != nil {
		return err
	}

	// Render the centralized routes file for the app (RegisterRoutes)
	routePath := "templates/app/routes_sample.tmpl"
	if req.IncludeExample {
		routePath = "templates/app/routes.tmpl"
	}
	routesOut := filepath.Join(tmp, "internal/app", "routes.go")
	routesData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Framework":  req.Framework,
	}
	if err := s.renderTemplate(routePath, routesOut, routesData); err != nil {
		return fmt.Errorf("failed to render routes template: %w", err)
	}

	// Render bootstrap that initializes repositories/usecases/handlers
	bootstrapPath := "templates/app/bootstrap_sample.tmpl"
	if req.IncludeExample {
		bootstrapPath = "templates/app/bootstrap.tmpl"
	}
	bootstrapOut := filepath.Join(tmp, "internal/app", "bootstrap.go")
	if err := s.renderTemplate(bootstrapPath, bootstrapOut, data); err != nil {
		return fmt.Errorf("failed to render bootstrap template: %w", err)
	}

	return nil
}

func (s *GeneratorService) renderJobsLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, "internal/jobs", "example_job.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate("templates/jobs/example_job.tmpl", outPath, data)
}

func loadManifest(path string) (*models.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m models.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *GeneratorService) GetManifest() *models.Manifest {
	return s.manifest
}

type GenerateRequest = models.GenerateRequest

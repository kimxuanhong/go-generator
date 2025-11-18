package service

import (
	"os"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
	"github.com/xhkzeroone/go-generator/internal/models"
)

type GeneratorService struct {
	manifest *models.Manifest
}

func NewGeneratorService(manifestPath string) (*GeneratorService, error) {
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		return nil, errors.ErrConfig("Failed to load manifest", err).
			WithContext("manifest_path", manifestPath)
	}
	return &GeneratorService{manifest: manifest}, nil
}

func (s *GeneratorService) GenerateProject(req *GenerateRequest) ([]byte, error) {
	// Validate framework exists
	if _, ok := s.manifest.Frameworks[req.Framework]; !ok {
		return nil, errors.ErrNotFound("framework").
			WithContext("framework", req.Framework)
	}

	// Validate all libs exist
	for _, lib := range req.Libs {
		if _, ok := s.manifest.Libs[lib]; !ok {
			return nil, errors.ErrNotFound("library").
				WithContext("library", lib)
		}
	}

	// Create temp directory
	tmp, err := os.MkdirTemp("", constants.TempDirPrefix+req.ProjectName+"-*")
	if err != nil {
		return nil, errors.ErrFileSystem("Failed to create temporary directory", err).
			WithContext("project_name", req.ProjectName)
	}
	defer os.RemoveAll(tmp)

	// Create base structure
	if err := s.createBaseStructure(tmp, req.IncludeExample, req.Framework); err != nil {
		return nil, errors.ErrFileSystem("Failed to create project structure", err).
			WithContext("project_name", req.ProjectName)
	}

	// Collect dependencies
	allImports := s.collectDependencies(req)

	// Get framework definition
	fdef := s.manifest.Frameworks[req.Framework]

	// Render framework templates
	if err := s.renderFrameworkTemplates(tmp, req, fdef); err != nil {
		return nil, errors.ErrTemplate("Failed to render framework templates", err).
			WithContext("framework", req.Framework)
	}

	// Render middleware templates (always included for logging, tracing, rate limiting)
	if err := s.renderMiddlewareTemplates(tmp, req); err != nil {
		return nil, errors.ErrTemplate("Failed to render middleware templates", err).
			WithContext("framework", req.Framework)
	}

	// Render library templates and merge configs
	includes := make(map[string]bool)
	mergedConfig := make(map[string]interface{})

	// Add default log config section (always included)
	mergedConfig["log"] = map[string]interface{}{
		"level": "info",
	}

	// Merge framework config section first
	if fdef.ConfigSection != "" {
		if err := s.mergeConfigSection(fdef.ConfigSection, mergedConfig); err != nil {
			return nil, errors.ErrConfig("Failed to merge framework configuration", err).
				WithContext("framework", req.Framework)
		}
	}

	// Merge library config sections
	for _, lib := range req.Libs {
		ldef := s.manifest.Libs[lib]
		includes[lib] = true

		if err := s.renderLibTemplates(tmp, req, lib, ldef); err != nil {
			return nil, errors.ErrTemplate("Failed to render library templates", err).
				WithContext("library", lib)
		}

		// Merge config section
		if ldef.ConfigSection != "" {
			if err := s.mergeConfigSection(ldef.ConfigSection, mergedConfig); err != nil {
				return nil, errors.ErrConfig("Failed to merge library configuration", err).
					WithContext("library", lib)
			}
		}
	}

	// Write config file
	if err := s.writeConfigFile(tmp, mergedConfig); err != nil {
		return nil, errors.ErrFileSystem("Failed to write configuration file", err)
	}

	// Render Clean Architecture layers only if IncludeExample is true
	if req.IncludeExample {
		if err := s.renderDomainLayer(tmp, req, includes); err != nil {
			return nil, errors.ErrTemplate("Failed to render domain layer", err)
		}

		// Render errors package (always included with examples)
		if err := s.renderErrorsLayer(tmp, req); err != nil {
			return nil, errors.ErrTemplate("Failed to render errors layer", err)
		}

		// Render database models (infrastructure layer)
		if err := s.renderModelsLayer(tmp, req, includes); err != nil {
			return nil, errors.ErrTemplate("Failed to render models layer", err)
		}

		if err := s.renderRepositoryLayer(tmp, req, includes); err != nil {
			return nil, errors.ErrTemplate("Failed to render repository layer", err)
		}

		if err := s.renderUsecaseLayer(tmp, req, includes); err != nil {
			return nil, errors.ErrTemplate("Failed to render usecase layer", err)
		}

		if err := s.renderHandlerLayer(tmp, req, includes); err != nil {
			return nil, errors.ErrTemplate("Failed to render handler layer", err)
		}

		// Jobs layer (only if cron is included)
		if includes["cron"] {
			if err := s.renderJobsLayer(tmp, req, includes); err != nil {
				return nil, errors.ErrTemplate("Failed to render jobs layer", err)
			}
		}

		// Consumers layer (if RabbitMQ, Kafka, or ActiveMQ is included)
		if includes["rabbitmq"] || includes["kafka"] || includes["activemq"] {
			if err := s.renderConsumersLayer(tmp, req, includes); err != nil {
				return nil, errors.ErrTemplate("Failed to render consumers layer", err)
			}
		}
	}

	// App server (always render, but with or without example routes)
	if err := s.renderAppServer(tmp, req, includes); err != nil {
		return nil, errors.ErrTemplate("Failed to render app server", err)
	}

	// Write main.go
	if err := s.renderMainFile(tmp, req, includes); err != nil {
		return nil, errors.ErrTemplate("Failed to render main file", err)
	}

	// Write deps package
	if err := s.renderDepsPackage(tmp, req, includes); err != nil {
		return nil, errors.ErrTemplate("Failed to render dependencies package", err)
	}

	// Render Swagger docs stub
	if err := s.renderDocsStub(tmp, req); err != nil {
		return nil, errors.ErrTemplate("Failed to render documentation", err)
	}

	// Write go.mod with all dependencies
	if err := s.renderGoMod(tmp, req, allImports); err != nil {
		return nil, errors.ErrTemplate("Failed to render go.mod file", err)
	}

	// Render project files (Dockerfile, .gitignore, .env.example, README.md)
	if err := s.renderProjectFiles(tmp, req, includes); err != nil {
		return nil, errors.ErrTemplate("Failed to render project files", err)
	}

	// Create zip file
	zipData, err := s.createZip(tmp)
	if err != nil {
		return nil, errors.ErrFileSystem("Failed to create project archive", err)
	}

	return zipData, nil
}

type GenerateRequest = models.GenerateRequest

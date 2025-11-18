package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
)

// renderDepsPackage renders the deps package using metadata
func (s *GeneratorService) renderDepsPackage(tmp string, req *GenerateRequest, includes map[string]bool) error {
	// Load metadata
	depsMeta, err := s.loadDepsMetadata()
	if err != nil {
		return err
	}

	// Replace {{.ModuleName}} in imports and remove duplicates across all depsMeta
	globalSeen := make(map[string]struct{})

	for key, meta := range depsMeta {
		uniqueImports := make([]string, 0, len(meta.Imports))
		for _, imp := range meta.Imports {
			imp = strings.ReplaceAll(imp, constants.ModuleNamePlaceholder, req.ModuleName)
			if _, exists := globalSeen[imp]; !exists {
				globalSeen[imp] = struct{}{}
				uniqueImports = append(uniqueImports, imp)
			}
		}
		meta.Imports = uniqueImports
		depsMeta[key] = meta
	}

	// Create deps directory
	depsDir := filepath.Join(tmp, constants.DirInternalDeps)
	if err := os.MkdirAll(depsDir, constants.DirPerm); err != nil {
		return errors.ErrFileSystem("Failed to create deps directory", err).
			WithContext("directory", depsDir)
	}

	// Render main deps.go
	depsPath := filepath.Join(depsDir, "deps.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
		"DepsMeta":   depsMeta,
	}

	if err := s.renderTemplate(constants.TemplateDeps, depsPath, data); err != nil {
		return err
	}

	// Render helper files for each included dependency
	for key, include := range includes {
		if !include {
			continue
		}

		meta, exists := depsMeta[key]
		if !exists {
			continue
		}

		// Render each helper file
		for _, helperFile := range meta.HelperFiles {
			templatePath := filepath.Join(constants.TemplateDepsDir, filepath.Base(helperFile))
			outputPath := filepath.Join(depsDir, filepath.Base(helperFile))

			// Remove .tmpl extension from output
			outputPath = strings.TrimSuffix(outputPath, constants.TemplateExtension) + constants.GoFileExtension

			helperData := map[string]interface{}{
				"ModuleName": req.ModuleName,
				"Key":        key,
				"Includes":   includes,
			}

			if err := s.renderTemplate(templatePath, outputPath, helperData); err != nil {
				return err
			}
		}
	}

	// Load config metadata
	configMeta, err := s.loadConfigMetadata()
	if err != nil {
		return err
	}

	// Replace {{.ModuleName}} with actual module name
	configMeta = applyModuleName(configMeta, req.ModuleName)

	// Render config.go
	configPath := filepath.Join(tmp, constants.DirInternalDeps, "config.go")
	configData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
		"Framework":  req.Framework,
		"Meta":       configMeta,
	}
	return s.renderTemplate(constants.TemplateConfig, configPath, configData)
}

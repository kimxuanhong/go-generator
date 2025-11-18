package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
	"github.com/xhkzeroone/go-generator/internal/models"
)

const (
	// ManifestVersion is the current supported manifest version
	ManifestVersion = "1.0.0"
	// MinManifestVersion is the minimum supported version
	MinManifestVersion = "1.0.0"
)

// validateManifest validates the manifest structure and content
func validateManifest(m *models.Manifest) error {
	// Validate version
	if m.Version == "" {
		return errors.ErrConfig("Manifest version is required", nil)
	}

	// Check version compatibility (simple string comparison for now)
	// In production, you might want to use semantic versioning library
	if m.Version != ManifestVersion && m.Version < MinManifestVersion {
		return errors.ErrConfig(
			fmt.Sprintf("Manifest version %s is not supported. Supported version is %s", m.Version, ManifestVersion),
			nil,
		).WithContext("manifest_version", m.Version).
			WithContext("supported_version", ManifestVersion)
	}

	// Validate frameworks
	if len(m.Frameworks) == 0 {
		return errors.ErrConfig("Manifest must contain at least one framework", nil)
	}

	for name, framework := range m.Frameworks {
		if err := validateFramework(name, framework); err != nil {
			return errors.ErrConfig(fmt.Sprintf("Invalid framework '%s': %v", name, err), nil)
		}
	}

	// Validate libraries
	for name, lib := range m.Libs {
		if err := validateLibrary(name, lib); err != nil {
			return errors.ErrConfig(fmt.Sprintf("Invalid library '%s': %v", name, err), nil)
		}
	}

	return nil
}

// validateFramework validates a framework definition
func validateFramework(name string, f models.FrameworkDef) error {
	if name == "" {
		return fmt.Errorf("framework name cannot be empty")
	}

	if len(f.Templates) == 0 {
		return fmt.Errorf("framework must have at least one template")
	}

	// Validate template paths exist
	for _, templatePath := range f.Templates {
		if !strings.HasSuffix(templatePath, constants.TemplateExtension) {
			return fmt.Errorf("template path must end with %s: %s", constants.TemplateExtension, templatePath)
		}
		if !strings.HasPrefix(templatePath, constants.TemplateDir+"/") {
			return fmt.Errorf("template path must be under %s/: %s", constants.TemplateDir, templatePath)
		}
		if _, err := os.Stat(templatePath); err != nil {
			return fmt.Errorf("template file does not exist: %s", templatePath)
		}
	}

	// Validate config section if provided
	if f.ConfigSection != "" {
		if !strings.HasSuffix(f.ConfigSection, ".json") {
			return fmt.Errorf("config section must be a JSON file: %s", f.ConfigSection)
		}
		if !strings.HasPrefix(f.ConfigSection, constants.TemplateDir+"/") {
			return fmt.Errorf("config section path must be under %s/: %s", constants.TemplateDir, f.ConfigSection)
		}
		if _, err := os.Stat(f.ConfigSection); err != nil {
			return fmt.Errorf("config section file does not exist: %s", f.ConfigSection)
		}
	}

	// Validate imports
	if len(f.Imports) == 0 {
		return fmt.Errorf("framework must have at least one import")
	}

	return nil
}

// validateLibrary validates a library definition
func validateLibrary(name string, l models.LibDef) error {
	if name == "" {
		return fmt.Errorf("library name cannot be empty")
	}

	if len(l.Templates) == 0 {
		return fmt.Errorf("library must have at least one template")
	}

	// Validate template paths exist
	for _, templatePath := range l.Templates {
		if !strings.HasSuffix(templatePath, constants.TemplateExtension) {
			return fmt.Errorf("template path must end with %s: %s", constants.TemplateExtension, templatePath)
		}
		if !strings.HasPrefix(templatePath, constants.TemplateDir+"/") {
			return fmt.Errorf("template path must be under %s/: %s", constants.TemplateDir, templatePath)
		}
		// Check if template exists
		if _, err := os.Stat(templatePath); err != nil {
			return fmt.Errorf("template file does not exist: %s", templatePath)
		}
	}

	// Validate config section if provided
	if l.ConfigSection != "" {
		if !strings.HasSuffix(l.ConfigSection, ".json") {
			return fmt.Errorf("config section must be a JSON file: %s", l.ConfigSection)
		}
		if !strings.HasPrefix(l.ConfigSection, constants.TemplateDir+"/") {
			return fmt.Errorf("config section path must be under %s/: %s", constants.TemplateDir, l.ConfigSection)
		}
		// Check if config section exists
		if _, err := os.Stat(l.ConfigSection); err != nil {
			return fmt.Errorf("config section file does not exist: %s", l.ConfigSection)
		}
	}

	// Validate imports
	if len(l.Imports) == 0 {
		return fmt.Errorf("library must have at least one import")
	}

	// Validate category if provided
	validCategories := map[string]bool{
		"database":      true,
		"caching":       true,
		"messaging":     true,
		"utilities":     true,
		"observability": true,
		"other":         true,
	}
	if l.Category != "" && !validCategories[l.Category] {
		return fmt.Errorf("invalid category: %s. Valid categories: database, caching, messaging, utilities, observability, other", l.Category)
	}

	return nil
}

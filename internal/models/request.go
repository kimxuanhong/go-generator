package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xhkzeroone/go-generator/internal/constants"
)

type GenerateRequest struct {
	ProjectName    string   `json:"projectName"`
	ModuleName     string   `json:"moduleName"`
	Framework      string   `json:"framework"`
	Architecture   string   `json:"architecture"`
	Libs           []string `json:"libs"`
	IncludeExample bool     `json:"includeExample,omitempty"` // Optional: include example code (User entity, usecase, handler)
}

func (r *GenerateRequest) Validate() error {
	if err := r.validateProjectName(); err != nil {
		return err
	}
	if err := r.validateModuleName(); err != nil {
		return err
	}
	if r.Framework == "" {
		return fmt.Errorf("framework is required")
	}
	return nil
}

func (r *GenerateRequest) validateProjectName() error {
	if r.ProjectName == "" {
		return fmt.Errorf("projectName is required")
	}

	// Check length
	if len(r.ProjectName) < constants.MinProjectNameLength {
		return fmt.Errorf("projectName must be at least %d characters", constants.MinProjectNameLength)
	}
	if len(r.ProjectName) > constants.MaxProjectNameLength {
		return fmt.Errorf("projectName must be at most %d characters", constants.MaxProjectNameLength)
	}

	// Check for path traversal attempts
	if strings.Contains(r.ProjectName, "..") || strings.Contains(r.ProjectName, "/") || strings.Contains(r.ProjectName, "\\") {
		return fmt.Errorf("projectName contains invalid characters (path traversal detected)")
	}

	// Check for null bytes
	if strings.Contains(r.ProjectName, "\x00") {
		return fmt.Errorf("projectName contains invalid characters")
	}

	// Validate pattern: lowercase alphanumeric with hyphens
	matched, err := regexp.MatchString(constants.ProjectNamePattern, r.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to validate projectName: %w", err)
	}
	if !matched {
		return fmt.Errorf("projectName must be lowercase alphanumeric with hyphens only (e.g., my-project)")
	}

	return nil
}

func (r *GenerateRequest) validateModuleName() error {
	if r.ModuleName == "" {
		return fmt.Errorf("moduleName is required")
	}

	// Check length
	if len(r.ModuleName) < constants.MinModuleNameLength {
		return fmt.Errorf("moduleName must be at least %d characters", constants.MinModuleNameLength)
	}
	if len(r.ModuleName) > constants.MaxModuleNameLength {
		return fmt.Errorf("moduleName must be at most %d characters", constants.MaxModuleNameLength)
	}

	// Check for path traversal attempts
	if strings.Contains(r.ModuleName, "..") {
		return fmt.Errorf("moduleName contains invalid characters (path traversal detected)")
	}

	// Check for null bytes
	if strings.Contains(r.ModuleName, "\x00") {
		return fmt.Errorf("moduleName contains invalid characters")
	}

	// Basic validation: should look like a Go module path
	// Allow: github.com/user/repo, example.com/pkg, etc.
	matched, err := regexp.MatchString(constants.ModuleNamePattern, r.ModuleName)
	if err != nil {
		return fmt.Errorf("failed to validate moduleName: %w", err)
	}
	if !matched {
		return fmt.Errorf("moduleName must be a valid Go module path (e.g., github.com/user/project)")
	}

	return nil
}

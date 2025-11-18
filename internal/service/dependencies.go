package service

import (
	"sort"
	"strings"

	"github.com/xhkzeroone/go-generator/internal/constants"
)

// collectDependencies collects all module dependencies from framework and libraries
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
	moduleMap[constants.DepLogrus] = true
	moduleMap[constants.DepViper] = true
	moduleMap[constants.DepGoogleUUID] = true

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

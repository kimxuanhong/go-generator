package service

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
)

// DepMetadata represents metadata for a dependency
type DepMetadata struct {
	Imports     []string `json:"imports"`
	StructField string   `json:"struct_field"`
	InitLines   []string `json:"init_lines"`
	CloseLines  []string `json:"close_lines"`
	HelperFiles []string `json:"helper_files"`
}

// DepMeta represents metadata for config dependencies
type DepMeta struct {
	Imports     []string `json:"imports"`
	ConfigField string   `json:"config_field"`
}

// loadDepsMetadata loads the deps metadata from JSON file
func (s *GeneratorService) loadDepsMetadata() (map[string]DepMetadata, error) {
	data, err := os.ReadFile(constants.TemplateDepsMeta)
	if err != nil {
		return nil, errors.ErrFileSystem("Failed to read deps metadata file", err).
			WithContext("path", constants.TemplateDepsMeta)
	}

	var metadata map[string]DepMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, errors.ErrConfig("Failed to parse deps metadata", err).
			WithContext("path", constants.TemplateDepsMeta)
	}

	return metadata, nil
}

// loadConfigMetadata loads the config metadata from JSON file
func (s *GeneratorService) loadConfigMetadata() (map[string]DepMeta, error) {
	b, err := os.ReadFile(constants.TemplateConfigMeta)
	if err != nil {
		return nil, errors.ErrFileSystem("Failed to read config metadata file", err).
			WithContext("path", constants.TemplateConfigMeta)
	}

	var meta map[string]DepMeta
	if err := json.Unmarshal(b, &meta); err != nil {
		return nil, errors.ErrConfig("Failed to parse config metadata", err).
			WithContext("path", constants.TemplateConfigMeta)
	}

	return meta, nil
}

// applyModuleName replaces module name placeholder in metadata imports
func applyModuleName(meta map[string]DepMeta, moduleName string) map[string]DepMeta {
	res := make(map[string]DepMeta)
	for k, v := range meta {
		newImports := make([]string, len(v.Imports))
		for i, imp := range v.Imports {
			newImports[i] = strings.ReplaceAll(imp, constants.ModuleNamePlaceholder, moduleName)
		}
		v.Imports = newImports
		res[k] = v
	}
	return res
}

package service

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
)

// mergeConfigSection merges a config section from a JSON file into the merged config
func (s *GeneratorService) mergeConfigSection(configPath string, mergedConfig map[string]interface{}) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return errors.ErrFileSystem("Failed to read config section file", err).
			WithContext("config_path", configPath)
	}

	var section map[string]interface{}
	if err := json.Unmarshal(data, &section); err != nil {
		return errors.ErrConfig("Failed to parse config section", err).
			WithContext("config_path", configPath)
	}

	for k, v := range section {
		mergedConfig[k] = v
	}

	return nil
}

// writeConfigFile writes the merged configuration to a JSON file
func (s *GeneratorService) writeConfigFile(tmp string, config map[string]interface{}) error {
	cfgPath := filepath.Join(tmp, constants.DirConfig, constants.ConfigFileName)
	cfgFile, err := os.Create(cfgPath)
	if err != nil {
		return errors.ErrFileSystem("Failed to create config file", err).
			WithContext("config_path", cfgPath)
	}
	defer cfgFile.Close()

	enc := json.NewEncoder(cfgFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		return errors.ErrFileSystem("Failed to write config file", err).
			WithContext("config_path", cfgPath)
	}
	return nil
}

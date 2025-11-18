package service

import (
	"encoding/json"
	"os"

	"github.com/xhkzeroone/go-generator/internal/errors"
	"github.com/xhkzeroone/go-generator/internal/models"
)

// loadManifest loads and validates the manifest from a JSON file
func loadManifest(path string) (*models.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.ErrFileSystem("Failed to read manifest file", err).
			WithContext("path", path)
	}

	var m models.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.ErrConfig("Failed to parse manifest file", err).
			WithContext("path", path)
	}

	// Validate manifest structure
	if err := validateManifest(&m); err != nil {
		return nil, errors.ErrConfig("Manifest validation failed", err).
			WithContext("path", path).
			WithContext("version", m.Version)
	}

	return &m, nil
}

// GetManifest returns the loaded manifest
func (s *GeneratorService) GetManifest() *models.Manifest {
	return s.manifest
}

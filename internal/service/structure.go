package service

import (
	"os"
	"path/filepath"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
)

// createBaseStructure creates the base directory structure for the project
func (s *GeneratorService) createBaseStructure(tmp string, includeExample bool, framework string) error {
	dirs := []string{
		constants.DirCmd,
		constants.DirDocs,
		constants.DirInternalApp,
		constants.DirInternalInfra,
		constants.DirInternalDeps,
		constants.DirInternalMiddleware,
		constants.DirConfig,
	}

	// Only create example layers if IncludeExample is true
	if includeExample {
		exampleDirs := []string{
			constants.DirInternalDomain,
			constants.DirInternalErrors,
			constants.DirInternalUsecase,
			constants.DirInternalRepo,
			constants.DirInternalRepoModels,
			constants.DirInternalHandler,
			constants.DirInternalJobs,
		}
		dirs = append(dirs, exampleDirs...)
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(tmp, dir)
		if err := os.MkdirAll(dirPath, constants.DirPerm); err != nil {
			return errors.ErrFileSystem("Failed to create directory", err).
				WithContext("directory", dirPath)
		}
	}
	return nil
}

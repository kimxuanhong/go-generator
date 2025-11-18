package service

import (
	"path/filepath"

	"github.com/xhkzeroone/go-generator/internal/constants"
)

// renderProjectFiles renders additional project files (Dockerfile, .gitignore, etc.)
func (s *GeneratorService) renderProjectFiles(tmp string, req *GenerateRequest, includes map[string]bool) error {
	// Render Dockerfile
	dockerData := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
		"BinaryName":  req.ProjectName,
		"Port":        constants.DefaultPortNum,
		"GoVersion":   constants.DefaultGoVersion,
	}
	if err := s.renderTemplate(constants.TemplateDockerfile, filepath.Join(tmp, constants.DockerfileName), dockerData); err != nil {
		return err
	}

	// Render .gitignore
	gitignoreData := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
	}
	if err := s.renderTemplate(constants.TemplateGitignore, filepath.Join(tmp, constants.GitignoreFileName), gitignoreData); err != nil {
		return err
	}

	// Render .env.example
	envExampleData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	if err := s.renderTemplate(constants.TemplateEnvExample, filepath.Join(tmp, constants.EnvExampleFileName), envExampleData); err != nil {
		return err
	}

	// Render README.md
	readmeData := map[string]interface{}{
		"ProjectName":    req.ProjectName,
		"ModuleName":     req.ModuleName,
		"Framework":      req.Framework,
		"IncludeExample": req.IncludeExample,
		"Includes":       includes,
	}
	if err := s.renderTemplate(constants.TemplateReadme, filepath.Join(tmp, constants.ReadmeFileName), readmeData); err != nil {
		return err
	}

	return nil
}

package service

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
	"github.com/xhkzeroone/go-generator/internal/models"
)

// renderTemplate renders a template file to a destination path
func (s *GeneratorService) renderTemplate(tmplPath, dstPath string, data interface{}) error {
	tpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.ErrTemplate("Failed to parse template", err).
			WithContext("template_path", tmplPath)
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), constants.DirPerm); err != nil {
		return errors.ErrFileSystem("Failed to create directory", err).
			WithContext("path", filepath.Dir(dstPath))
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return errors.ErrFileSystem("Failed to create file", err).
			WithContext("path", dstPath)
	}
	defer f.Close()

	if err := tpl.Execute(f, data); err != nil {
		return errors.ErrTemplate("Failed to execute template", err).
			WithContext("template_path", tmplPath).
			WithContext("output_path", dstPath)
	}

	return nil
}

// renderFrameworkTemplates renders framework-specific templates
func (s *GeneratorService) renderFrameworkTemplates(tmp string, req *GenerateRequest, fdef models.FrameworkDef) error {
	for _, t := range fdef.Templates {
		baseName := filepath.Base(strings.TrimSuffix(t, constants.TemplateExtension))
		outPath := filepath.Join(tmp, constants.DirInternalApp, baseName+constants.GoFileExtension)

		data := map[string]interface{}{
			"ModuleName":  req.ModuleName,
			"ProjectName": req.ProjectName,
			"Framework":   req.Framework,
		}
		if err := s.renderTemplate(t, outPath, data); err != nil {
			return err
		}
	}
	return nil
}

// renderLibTemplates renders library-specific templates
func (s *GeneratorService) renderLibTemplates(tmp string, req *GenerateRequest, lib string, ldef models.LibDef) error {
	for _, t := range ldef.Templates {
		baseName := filepath.Base(strings.TrimSuffix(t, constants.TemplateExtension))
		outPath := filepath.Join(tmp, constants.DirInternalInfra, lib, baseName+constants.GoFileExtension)

		data := map[string]interface{}{
			"ModuleName":  req.ModuleName,
			"ProjectName": req.ProjectName,
			"Lib":         lib,
		}
		if err := s.renderTemplate(t, outPath, data); err != nil {
			return err
		}
	}
	return nil
}

// renderMainFile renders the main.go file
func (s *GeneratorService) renderMainFile(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, constants.DirCmd, "main.go")
	data := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
		"Framework":   req.Framework,
		"Includes":    includes,
	}
	return s.renderTemplate(constants.TemplateMain, outPath, data)
}

// renderDocsStub renders the Swagger docs stub
func (s *GeneratorService) renderDocsStub(tmp string, req *GenerateRequest) error {
	outPath := filepath.Join(tmp, constants.DirDocs, "docs.go")
	data := map[string]interface{}{
		"ModuleName":  req.ModuleName,
		"ProjectName": req.ProjectName,
	}
	return s.renderTemplate(constants.TemplateDocs, outPath, data)
}

// renderGoMod renders the go.mod file
func (s *GeneratorService) renderGoMod(tmp string, req *GenerateRequest, imports []string) error {
	outPath := filepath.Join(tmp, constants.GoModFileName)
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Imports":    imports,
	}
	return s.renderTemplate(constants.TemplateGoMod, outPath, data)
}

// renderMiddlewareTemplates renders middleware templates for the selected framework
func (s *GeneratorService) renderMiddlewareTemplates(tmp string, req *GenerateRequest) error {
	// Define middleware template files based on framework
	middlewareDir := filepath.Join("templates", "middleware", req.Framework)
	middlewareFiles := []string{"logging.tmpl", "tracing.tmpl", "ratelimit.tmpl"}

	for _, middlewareFile := range middlewareFiles {
		tmplPath := filepath.Join(middlewareDir, middlewareFile)
		baseName := filepath.Base(strings.TrimSuffix(middlewareFile, constants.TemplateExtension))
		outPath := filepath.Join(tmp, constants.DirInternalMiddleware, baseName+constants.GoFileExtension)

		data := map[string]interface{}{
			"ModuleName":  req.ModuleName,
			"ProjectName": req.ProjectName,
			"Framework":   req.Framework,
		}

		if err := s.renderTemplate(tmplPath, outPath, data); err != nil {
			return err
		}
	}
	return nil
}

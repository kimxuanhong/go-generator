package models

import "fmt"

type GenerateRequest struct {
	ProjectName    string   `json:"projectName"`
	ModuleName     string   `json:"moduleName"`
	Framework      string   `json:"framework"`
	Architecture   string   `json:"architecture"`
	Libs           []string `json:"libs"`
	IncludeExample bool     `json:"includeExample,omitempty"` // Optional: include example code (User entity, usecase, handler)
}

func (r *GenerateRequest) Validate() error {
	if r.ProjectName == "" {
		return fmt.Errorf("projectName is required")
	}
	if r.ModuleName == "" {
		return fmt.Errorf("moduleName is required")
	}
	if r.Framework == "" {
		return fmt.Errorf("framework is required")
	}
	return nil
}

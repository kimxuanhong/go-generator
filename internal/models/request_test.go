package models

import (
	"testing"
)

func TestGenerateRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: GenerateRequest{
				ProjectName:    "my-project",
				ModuleName:     "github.com/user/my-project",
				Framework:      "gin",
				Libs:           []string{"redis", "postgres"},
				IncludeExample: true,
			},
			wantErr: false,
		},
		{
			name: "missing project name",
			req: GenerateRequest{
				ModuleName: "github.com/user/my-project",
				Framework:  "gin",
			},
			wantErr: true,
			errMsg:  "projectName is required",
		},
		{
			name: "missing module name",
			req: GenerateRequest{
				ProjectName: "my-project",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "moduleName is required",
		},
		{
			name: "missing framework",
			req: GenerateRequest{
				ProjectName: "my-project",
				ModuleName:  "github.com/user/my-project",
			},
			wantErr: true,
			errMsg:  "framework is required",
		},
		{
			name: "project name with path traversal",
			req: GenerateRequest{
				ProjectName: "../etc/passwd",
				ModuleName:  "github.com/user/my-project",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "path traversal",
		},
		{
			name: "project name with uppercase",
			req: GenerateRequest{
				ProjectName: "My-Project",
				ModuleName:  "github.com/user/my-project",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "lowercase alphanumeric",
		},
		{
			name: "project name with special characters",
			req: GenerateRequest{
				ProjectName: "my_project@123",
				ModuleName:  "github.com/user/my-project",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "lowercase alphanumeric",
		},
		{
			name: "module name with path traversal",
			req: GenerateRequest{
				ProjectName: "my-project",
				ModuleName:  "../../etc/passwd",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "path traversal",
		},
		{
			name: "invalid module name format",
			req: GenerateRequest{
				ProjectName: "my-project",
				ModuleName:  "invalid-module-name",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "valid Go module path",
		},
		{
			name: "project name too short",
			req: GenerateRequest{
				ProjectName: "",
				ModuleName:  "github.com/user/my-project",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "projectName is required",
		},
		{
			name: "project name too long",
			req: GenerateRequest{
				ProjectName: "a" + string(make([]byte, 51)),
				ModuleName:  "github.com/user/my-project",
				Framework:   "gin",
			},
			wantErr: true,
			errMsg:  "at most",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("Validate() expected error message containing %q, got %v", tt.errMsg, err)
				} else if err != nil && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error message = %q, want containing %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestGenerateRequest_validateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateRequest
		wantErr bool
	}{
		{
			name: "valid project name",
			req: GenerateRequest{
				ProjectName: "my-awesome-project",
			},
			wantErr: false,
		},
		{
			name: "project name with numbers",
			req: GenerateRequest{
				ProjectName: "project123",
			},
			wantErr: false,
		},
		{
			name: "project name with null byte",
			req: GenerateRequest{
				ProjectName: "project\x00name",
			},
			wantErr: true,
		},
		{
			name: "project name with slash",
			req: GenerateRequest{
				ProjectName: "project/name",
			},
			wantErr: true,
		},
		{
			name: "project name with backslash",
			req: GenerateRequest{
				ProjectName: "project\\name",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.validateProjectName()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateRequest_validateModuleName(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateRequest
		wantErr bool
	}{
		{
			name: "valid module name",
			req: GenerateRequest{
				ModuleName: "github.com/user/repo",
			},
			wantErr: false,
		},
		{
			name: "valid module name with subpath",
			req: GenerateRequest{
				ModuleName: "github.com/user/repo/pkg/sub",
			},
			wantErr: false,
		},
		{
			name: "valid module name with version",
			req: GenerateRequest{
				ModuleName: "example.com/v2",
			},
			wantErr: false,
		},
		{
			name: "module name with path traversal",
			req: GenerateRequest{
				ModuleName: "github.com/../etc/passwd",
			},
			wantErr: true,
		},
		{
			name: "module name with null byte",
			req: GenerateRequest{
				ModuleName: "github.com/user\x00repo",
			},
			wantErr: true,
		},
		{
			name: "invalid module name format",
			req: GenerateRequest{
				ModuleName: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.validateModuleName()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateModuleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

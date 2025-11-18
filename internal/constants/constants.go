package constants

const (
	// Default values
	DefaultPort         = ":8080"
	DefaultManifestPath = "manifest.json"
	DefaultTempPrefix   = "gen-"

	// HTTP methods
	MethodGET  = "GET"
	MethodPOST = "POST"

	// Content types
	ContentTypeJSON = "application/json"
	ContentTypeZip  = "application/zip"

	// Headers
	HeaderContentType        = "Content-Type"
	HeaderContentDisposition = "Content-Disposition"
	HeaderCacheControl       = "Cache-Control"
	HeaderPragma             = "Pragma"
	HeaderExpires            = "Expires"

	// Cache control values
	NoCache = "no-cache, no-store, must-revalidate"

	// Error messages
	ErrMethodNotAllowed    = "Method not allowed"
	ErrInvalidRequestBody  = "Invalid request body"
	ErrValidationFailed    = "Validation error"
	ErrGenerationFailed    = "Failed to generate project"
	ErrEncodingFailed      = "Failed to encode response"
	ErrInternalServerError = "Internal server error"

	// Validation patterns
	ProjectNamePattern = `^[a-z0-9-]+$`
	ModuleNamePattern  = `^[a-z0-9][a-z0-9._-]*/[a-z0-9][a-z0-9._-]*(/[a-z0-9][a-z0-9._-]*)*$`

	// Validation constraints
	MinProjectNameLength = 1
	MaxProjectNameLength = 50
	MinModuleNameLength  = 3
	MaxModuleNameLength  = 200

	// Service constants
	TempDirPrefix         = "gen-"
	DirPerm               = 0755
	TemplateExtension     = ".tmpl"
	GoFileExtension       = ".go"
	ConfigFileName        = "config.json"
	GoModFileName         = "go.mod"
	ReadmeFileName        = "README.md"
	GitignoreFileName     = ".gitignore"
	EnvExampleFileName    = ".env.example"
	DockerfileName        = "Dockerfile"
	ModuleNamePlaceholder = "{{.ModuleName}}"

	// Directory paths
	DirCmd                = "cmd"
	DirDocs               = "docs"
	DirInternalApp        = "internal/app"
	DirInternalInfra      = "internal/infrastructure"
	DirInternalDeps       = "internal/deps"
	DirInternalMiddleware = "internal/middleware"
	DirConfig             = "config"
	DirInternalDomain     = "internal/domain"
	DirInternalErrors     = "internal/errors"
	DirInternalUsecase    = "internal/usecase"
	DirInternalRepo       = "internal/infrastructure/repository"
	DirInternalRepoModels = "internal/infrastructure/repository/models"
	DirInternalHandler    = "internal/adapter/handler"
	DirInternalJobs       = "internal/jobs"

	// Template paths
	TemplateDir             = "templates"
	TemplateDepsDir         = "templates/deps"
	TemplateDepsMeta        = "templates/deps/deps_meta.json"
	TemplateConfigMeta      = "templates/deps/config_meta.json"
	TemplateDockerfile      = "templates/Dockerfile.tmpl"
	TemplateGitignore       = "templates/gitignore.tmpl"
	TemplateEnvExample      = "templates/env_example.tmpl"
	TemplateReadme          = "templates/README.tmpl"
	TemplateMain            = "templates/cmd/main.tmpl"
	TemplateDocs            = "templates/docs/swagger.tmpl"
	TemplateGoMod           = "templates/go_mod.tmpl"
	TemplateDomainEntity    = "templates/domain/entity.tmpl"
	TemplateErrors          = "templates/errors/errors.tmpl"
	TemplateUserModel       = "templates/infrastructure/models/user_model.tmpl"
	TemplateUserRepo        = "templates/infrastructure/repository/user_repository.tmpl"
	TemplateCacheRepo       = "templates/infrastructure/repository/cache_repository.tmpl"
	TemplateUserUsecase     = "templates/usecase/user_usecase.tmpl"
	TemplateUserHandler     = "templates/handler/user_handler.tmpl"
	TemplateExampleJob      = "templates/jobs/example_job.tmpl"
	TemplateServerSimple    = "templates/app/server_simple.tmpl"
	TemplateServer          = "templates/app/server.tmpl"
	TemplateRoutesSample    = "templates/app/routes_sample.tmpl"
	TemplateRoutes          = "templates/app/routes.tmpl"
	TemplateBootstrapSample = "templates/app/bootstrap_sample.tmpl"
	TemplateBootstrap       = "templates/app/bootstrap.tmpl"
	TemplateDeps            = "templates/deps/deps.tmpl"
	TemplateConfig          = "templates/deps/config.tmpl"

	// Common dependencies
	DepLogrus     = "github.com/sirupsen/logrus"
	DepViper      = "github.com/spf13/viper"
	DepGoogleUUID = "github.com/google/uuid"

	// Default values for generated projects
	DefaultGoVersion = "1.20"
	DefaultPortNum   = 8080
)

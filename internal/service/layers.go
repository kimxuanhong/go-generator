package service

import (
	"path/filepath"

	"github.com/xhkzeroone/go-generator/internal/constants"
)

// renderDomainLayer renders the domain layer templates
func (s *GeneratorService) renderDomainLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, constants.DirInternalDomain, "entity.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate(constants.TemplateDomainEntity, outPath, data)
}

// renderErrorsLayer renders the errors package
func (s *GeneratorService) renderErrorsLayer(tmp string, req *GenerateRequest) error {
	outPath := filepath.Join(tmp, constants.DirInternalErrors, "errors.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
	}
	return s.renderTemplate(constants.TemplateErrors, outPath, data)
}

// renderModelsLayer renders the database models (infrastructure layer)
func (s *GeneratorService) renderModelsLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, constants.DirInternalRepoModels, "user_model.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate(constants.TemplateUserModel, outPath, data)
}

// renderRepositoryLayer renders the repository layer templates
func (s *GeneratorService) renderRepositoryLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	// Render user repository
	userRepoPath := filepath.Join(tmp, constants.DirInternalRepo, "user_repository.go")
	userRepoData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	if err := s.renderTemplate(constants.TemplateUserRepo, userRepoPath, userRepoData); err != nil {
		return err
	}

	// Render cache repository
	cacheRepoPath := filepath.Join(tmp, constants.DirInternalRepo, "cache_repository.go")
	cacheRepoData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate(constants.TemplateCacheRepo, cacheRepoPath, cacheRepoData)
}

// renderUsecaseLayer renders the usecase layer templates
func (s *GeneratorService) renderUsecaseLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, constants.DirInternalUsecase, "user_usecase.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}
	return s.renderTemplate(constants.TemplateUserUsecase, outPath, data)
}

// renderHandlerLayer renders the handler layer templates
func (s *GeneratorService) renderHandlerLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, constants.DirInternalHandler, "user_handler.go")
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Framework":  req.Framework,
		"Includes":   includes,
	}
	return s.renderTemplate(constants.TemplateUserHandler, outPath, data)
}

// renderJobsLayer renders the scheduled jobs layer templates (Input Adapter: Jobs)
func (s *GeneratorService) renderJobsLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}

	// Render example job (Adapter: Scheduled Jobs)
	jobPath := filepath.Join(tmp, constants.DirInternalJob, "example_job.go")
	return s.renderTemplate(constants.TemplateExampleJob, jobPath, data)
}

// renderConsumersLayer renders the message queue consumers layer templates (Input Adapter: Consumers)
func (s *GeneratorService) renderConsumersLayer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	data := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Includes":   includes,
	}

	// Render RabbitMQ consumer if RabbitMQ is included (Adapter: Message Consumer)
	if includes["rabbitmq"] {
		rabbitPath := filepath.Join(tmp, constants.DirInternalConsumer, "user_rabbitmq_consumer.go")
		if err := s.renderTemplate(constants.TemplateRabbitMQConsumer, rabbitPath, data); err != nil {
			return err
		}
	}

	// Render Kafka consumer if Kafka is included (Adapter: Message Consumer)
	if includes["kafka"] {
		kafkaPath := filepath.Join(tmp, constants.DirInternalConsumer, "user_kafka_consumer.go")
		if err := s.renderTemplate(constants.TemplateKafkaConsumer, kafkaPath, data); err != nil {
			return err
		}
	}

	// Render ActiveMQ consumer if ActiveMQ is included (Adapter: Message Consumer)
	if includes["activemq"] {
		activemqPath := filepath.Join(tmp, constants.DirInternalConsumer, "user_activemq_consumer.go")
		if err := s.renderTemplate(constants.TemplateActiveMQConsumer, activemqPath, data); err != nil {
			return err
		}
	}

	return nil
}

// renderAppServer renders the app server templates
func (s *GeneratorService) renderAppServer(tmp string, req *GenerateRequest, includes map[string]bool) error {
	outPath := filepath.Join(tmp, constants.DirInternalApp, "server.go")
	data := map[string]interface{}{
		"ModuleName":     req.ModuleName,
		"ProjectName":    req.ProjectName,
		"Framework":      req.Framework,
		"Includes":       includes,
		"IncludeExample": req.IncludeExample,
	}

	// Use example server template if IncludeExample is true, otherwise use simple server
	templatePath := constants.TemplateServerSimple
	if req.IncludeExample {
		templatePath = constants.TemplateServer
	}

	// Render server.go
	if err := s.renderTemplate(templatePath, outPath, data); err != nil {
		return err
	}

	// Render the centralized routes file for the app (RegisterRoutes)
	routePath := constants.TemplateRoutesSample
	if req.IncludeExample {
		routePath = constants.TemplateRoutes
	}
	routesOut := filepath.Join(tmp, constants.DirInternalApp, "routes.go")
	routesData := map[string]interface{}{
		"ModuleName": req.ModuleName,
		"Framework":  req.Framework,
	}
	if err := s.renderTemplate(routePath, routesOut, routesData); err != nil {
		return err
	}

	// Render bootstrap that initializes repositories/usecases/handlers
	bootstrapPath := constants.TemplateBootstrapSample
	if req.IncludeExample {
		bootstrapPath = constants.TemplateBootstrap
	}
	bootstrapOut := filepath.Join(tmp, constants.DirInternalApp, "bootstrap.go")
	if err := s.renderTemplate(bootstrapPath, bootstrapOut, data); err != nil {
		return err
	}

	return nil
}

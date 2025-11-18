package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
	"github.com/xhkzeroone/go-generator/internal/middleware"
	"github.com/xhkzeroone/go-generator/internal/service"
)

type GenerateHandler struct {
	*BaseHandler
	service *service.GeneratorService
}

func NewGenerateHandler(svc *service.GeneratorService, logger *logrus.Logger) *GenerateHandler {
	return &GenerateHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     svc,
	}
}

// HandleGenerate godoc
// @Summary Generate a Go project scaffold
// @Description Generates a Go project scaffold based on the provided configuration and returns it as a ZIP archive.
// @Tags generator
// @Accept json
// @Produce application/zip
// @Param request body service.GenerateRequest true "Generator configuration"
// @Success 200 {file} file "Generated project archive"
// @Failure 400 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /generate [post]
func (h *GenerateHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(w)

	if !h.validateMethod(r, constants.MethodPOST) {
		h.writeErrorWithID(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed, requestID)
		return
	}

	var req service.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err,
		}).Warn("Invalid request body")
		h.writeErrorWithID(w, constants.ErrInvalidRequestBody, http.StatusBadRequest, requestID)
		return
	}

	if err := req.Validate(); err != nil {
		appErr := errors.ErrValidation(err.Error(), err).WithContext("project_name", req.ProjectName).
			WithContext("module_name", req.ModuleName).WithContext("framework", req.Framework)
		h.handleAppError(w, r, appErr)
		return
	}

	// Get trace context
	traceCtx := middleware.GetTraceContext(r.Context())

	// Start span for project generation
	ctx, spanID, finishSpan := middleware.StartSpan(r.Context(), "generate_project")
	defer finishSpan()

	logFields := logrus.Fields{
		"request_id":      requestID,
		"project_name":    req.ProjectName,
		"module_name":     req.ModuleName,
		"framework":       req.Framework,
		"libs":            req.Libs,
		"include_example": req.IncludeExample,
	}
	if traceCtx.TraceID != "" {
		logFields["trace_id"] = traceCtx.TraceID
	}
	if spanID != "" {
		logFields["span_id"] = spanID
	}
	h.logger.WithFields(logFields).Info("Generating project")

	// Record generation start time
	startTime := time.Now()

	zipData, err := h.service.GenerateProject(&req)

	// Record metrics
	duration := time.Since(startTime)
	middleware.RecordProjectGeneration(req.Framework, duration, int64(len(zipData)), err == nil)
	if err != nil {
		// Check if it's already an AppError
		if appErr, ok := err.(*errors.AppError); ok {
			appErr.WithContext("request_id", requestID).
				WithContext("project_name", req.ProjectName).
				WithContext("module_name", req.ModuleName)
			h.handleAppError(w, r, appErr)
			return
		}

		// Wrap generic error
		appErr := errors.ErrGeneration(constants.ErrGenerationFailed, err).
			WithContext("request_id", requestID).
			WithContext("project_name", req.ProjectName).
			WithContext("module_name", req.ModuleName)
		h.handleAppError(w, r, appErr)
		return
	}

	w.Header().Set(constants.HeaderContentType, constants.ContentTypeZip)
	w.Header().Set(constants.HeaderContentDisposition, "attachment; filename="+req.ProjectName+".zip")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(zipData); err != nil {
		h.logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err,
		}).Error("Error writing zip data")
		return
	}

	logFields = logrus.Fields{
		"request_id":   requestID,
		"project_name": req.ProjectName,
		"size_bytes":   len(zipData),
		"duration_ms":  duration.Milliseconds(),
	}
	if traceCtx.TraceID != "" {
		logFields["trace_id"] = traceCtx.TraceID
	}
	if spanID != "" {
		logFields["span_id"] = spanID
	}
	h.logger.WithFields(logFields).Info("Project generated successfully")

	_ = ctx // Use context for future enhancements
}

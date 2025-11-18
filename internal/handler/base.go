package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/errors"
	"github.com/xhkzeroone/go-generator/internal/middleware"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	logger *logrus.Logger
}

// NewBaseHandler creates a new base handler with logger
func NewBaseHandler(logger *logrus.Logger) *BaseHandler {
	return &BaseHandler{logger: logger}
}

// writeError writes an error response in JSON format
func (h *BaseHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	requestID := middleware.GetRequestID(w)
	h.writeErrorWithID(w, message, statusCode, requestID)
}

// writeErrorWithID writes an error response with request ID
func (h *BaseHandler) writeErrorWithID(w http.ResponseWriter, message string, statusCode int, requestID string) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:     message,
		RequestID: requestID,
	})
}

// handleAppError handles AppError and writes appropriate response
func (h *BaseHandler) handleAppError(w http.ResponseWriter, r *http.Request, appErr *errors.AppError) {
	requestID := middleware.GetRequestID(w)

	// Log error with full context (internal details)
	logFields := logrus.Fields{
		"request_id": requestID,
		"error_code": appErr.Code,
		"message":    appErr.Message,
	}

	// Add context from error
	for k, v := range appErr.Context {
		logFields[k] = v
	}

	// Add internal error if present (for logging only)
	if appErr.InternalError != nil {
		logFields["internal_error"] = appErr.InternalError.Error()
	}

	// Determine status code based on error code
	statusCode := http.StatusInternalServerError
	switch appErr.Code {
	case errors.ErrCodeValidation:
		statusCode = http.StatusBadRequest
	case errors.ErrCodeNotFound:
		statusCode = http.StatusNotFound
	case errors.ErrCodeGeneration, errors.ErrCodeConfig, errors.ErrCodeTemplate, errors.ErrCodeFileSystem:
		statusCode = http.StatusInternalServerError
	}

	// Log based on severity
	if statusCode >= 500 {
		h.logger.WithFields(logFields).Error("Application error occurred")
	} else {
		h.logger.WithFields(logFields).Warn("Application error occurred")
	}

	// Return user-friendly message (no internal details)
	h.writeErrorWithID(w, appErr.Message, statusCode, requestID)
}

// writeJSON writes a JSON response
func (h *BaseHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
		http.Error(w, constants.ErrEncodingFailed, http.StatusInternalServerError)
	}
}

// validateMethod validates that the request uses the expected HTTP method
func (h *BaseHandler) validateMethod(r *http.Request, expectedMethod string) bool {
	return r.Method == expectedMethod
}

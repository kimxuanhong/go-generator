package errors

import (
	"fmt"
)

// AppError represents an application error with context
type AppError struct {
	// Code is a machine-readable error code
	Code string
	// Message is a user-friendly error message
	Message string
	// InternalError is the underlying error (not exposed to users)
	InternalError error
	// Context provides additional context for logging
	Context map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.InternalError)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.InternalError
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewAppError creates a new application error
func NewAppError(code, message string, internalErr error) *AppError {
	return &AppError{
		Code:          code,
		Message:       message,
		InternalError: internalErr,
		Context:       make(map[string]interface{}),
	}
}

// Error codes
const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeNotFound   = "NOT_FOUND"
	ErrCodeInternal   = "INTERNAL_ERROR"
	ErrCodeGeneration = "GENERATION_ERROR"
	ErrCodeConfig     = "CONFIG_ERROR"
	ErrCodeTemplate   = "TEMPLATE_ERROR"
	ErrCodeFileSystem = "FILESYSTEM_ERROR"
)

// Common error constructors
func ErrValidation(message string, internalErr error) *AppError {
	return NewAppError(ErrCodeValidation, message, internalErr)
}

func ErrNotFound(resource string) *AppError {
	return NewAppError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource), nil)
}

func ErrInternal(message string, internalErr error) *AppError {
	return NewAppError(ErrCodeInternal, message, internalErr)
}

func ErrGeneration(message string, internalErr error) *AppError {
	return NewAppError(ErrCodeGeneration, message, internalErr)
}

func ErrConfig(message string, internalErr error) *AppError {
	return NewAppError(ErrCodeConfig, message, internalErr)
}

func ErrTemplate(message string, internalErr error) *AppError {
	return NewAppError(ErrCodeTemplate, message, internalErr)
}

func ErrFileSystem(message string, internalErr error) *AppError {
	return NewAppError(ErrCodeFileSystem, message, internalErr)
}

package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDContextKey is the context key for request ID
	RequestIDContextKey = "request_id"
)

// LoggingMiddleware logs HTTP requests with structured logging
func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate or get request ID
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Add request ID to response header
			w.Header().Set(RequestIDHeader, requestID)

			// Create response writer wrapper to capture status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Get trace context if available
			traceID := GetTraceID(r.Context())
			spanID := GetSpanID(r.Context())

			// Log request
			fields := logrus.Fields{
				"request_id":  requestID,
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			}
			if traceID != "" {
				fields["trace_id"] = traceID
			}
			if spanID != "" {
				fields["span_id"] = spanID
			}
			logger.WithFields(fields).Info("Request started")

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Log response
			responseFields := logrus.Fields{
				"request_id":  requestID,
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": rw.statusCode,
				"duration_ms": duration.Milliseconds(),
			}
			if traceID != "" {
				responseFields["trace_id"] = traceID
			}
			if spanID != "" {
				responseFields["span_id"] = spanID
			}

			// Log level based on status code
			if rw.statusCode >= 500 {
				logger.WithFields(responseFields).Error("Request completed with server error")
			} else if rw.statusCode >= 400 {
				logger.WithFields(responseFields).Warn("Request completed with client error")
			} else {
				logger.WithFields(responseFields).Info("Request completed")
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// GetRequestID extracts request ID from response headers (for logging in handlers)
func GetRequestID(w http.ResponseWriter) string {
	return w.Header().Get(RequestIDHeader)
}

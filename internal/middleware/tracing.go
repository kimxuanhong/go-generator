package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	// TraceIDHeader is the header name for trace ID
	TraceIDHeader = "X-Trace-ID"
	// SpanIDHeader is the header name for span ID
	SpanIDHeader = "X-Span-ID"
	// TraceIDContextKey is the context key for trace ID
	TraceIDContextKey = "trace_id"
	// SpanIDContextKey is the context key for span ID
	SpanIDContextKey = "span_id"
	// ParentSpanIDContextKey is the context key for parent span ID
	ParentSpanIDContextKey = "parent_span_id"
)

// TraceContext holds tracing information
type TraceContext struct {
	TraceID      string
	SpanID       string
	ParentSpanID string
}

// TracingMiddleware adds distributed tracing support
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get or generate trace ID
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Get parent span ID if present
		parentSpanID := r.Header.Get(SpanIDHeader)

		// Generate new span ID
		spanID := uuid.New().String()

		// Add trace IDs to response headers
		w.Header().Set(TraceIDHeader, traceID)
		w.Header().Set(SpanIDHeader, spanID)

		// Add to request context
		ctx := context.WithValue(r.Context(), TraceIDContextKey, traceID)
		ctx = context.WithValue(ctx, SpanIDContextKey, spanID)
		if parentSpanID != "" {
			ctx = context.WithValue(ctx, ParentSpanIDContextKey, parentSpanID)
		}

		// Create new request with context
		r = r.WithContext(ctx)

		// Process request
		next.ServeHTTP(w, r)
	})
}

// GetTraceID extracts trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDContextKey).(string); ok {
		return traceID
	}
	return ""
}

// GetSpanID extracts span ID from context
func GetSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(SpanIDContextKey).(string); ok {
		return spanID
	}
	return ""
}

// GetParentSpanID extracts parent span ID from context
func GetParentSpanID(ctx context.Context) string {
	if parentSpanID, ok := ctx.Value(ParentSpanIDContextKey).(string); ok {
		return parentSpanID
	}
	return ""
}

// GetTraceContext extracts full trace context from context
func GetTraceContext(ctx context.Context) TraceContext {
	return TraceContext{
		TraceID:      GetTraceID(ctx),
		SpanID:       GetSpanID(ctx),
		ParentSpanID: GetParentSpanID(ctx),
	}
}

// StartSpan creates a new span for tracing
func StartSpan(ctx context.Context, operation string) (context.Context, string, func()) {
	parentSpanID := GetSpanID(ctx)
	spanID := uuid.New().String()
	traceID := GetTraceID(ctx)
	if traceID == "" {
		traceID = uuid.New().String()
	}

	newCtx := context.WithValue(ctx, SpanIDContextKey, spanID)
	if parentSpanID != "" {
		newCtx = context.WithValue(newCtx, ParentSpanIDContextKey, parentSpanID)
	}
	newCtx = context.WithValue(newCtx, TraceIDContextKey, traceID)

	start := time.Now()

	finish := func() {
		duration := time.Since(start)
		// In a full tracing implementation, you would send this to a tracing backend
		// For now, we just have the infrastructure in place
		_ = duration
	}

	return newCtx, spanID, finish
}

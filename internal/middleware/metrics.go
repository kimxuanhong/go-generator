package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
		},
		[]string{"method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"method", "path", "status_code"},
	)

	// Business metrics
	projectGenerationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "project_generation_total",
			Help: "Total number of project generations",
		},
		[]string{"framework", "status"},
	)

	projectGenerationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "project_generation_duration_seconds",
			Help:    "Project generation duration in seconds",
			Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60, 120},
		},
		[]string{"framework"},
	)

	projectGenerationSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "project_generation_size_bytes",
			Help:    "Generated project size in bytes",
			Buckets: []float64{10000, 50000, 100000, 500000, 1000000, 5000000, 10000000},
		},
		[]string{"framework"},
	)
)

// MetricsMiddleware collects Prometheus metrics for HTTP requests
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get request size
		requestSize := float64(0)
		if r.ContentLength > 0 {
			requestSize = float64(r.ContentLength)
		}

		// Record request size
		path := sanitizePath(r.URL.Path)
		httpRequestSize.WithLabelValues(r.Method, path).Observe(requestSize)

		// Create response writer wrapper to capture status code and size
		rw := &metricsResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Record metrics
		statusCode := strconv.Itoa(rw.statusCode)
		httpRequestsTotal.WithLabelValues(r.Method, path, statusCode).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path, statusCode).Observe(duration)
		httpResponseSize.WithLabelValues(r.Method, path, statusCode).Observe(float64(rw.size))
	})
}

// metricsResponseWriter wraps http.ResponseWriter to capture status code and size
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// RecordProjectGeneration records metrics for project generation
func RecordProjectGeneration(framework string, duration time.Duration, size int64, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	projectGenerationTotal.WithLabelValues(framework, status).Inc()
	if success {
		projectGenerationDuration.WithLabelValues(framework).Observe(duration.Seconds())
		projectGenerationSize.WithLabelValues(framework).Observe(float64(size))
	}
}

// sanitizePath sanitizes the path for metrics (removes dynamic parts)
func sanitizePath(path string) string {
	// Replace common dynamic parts
	if path == "/" {
		return "/"
	}
	// Keep only the first segment for most paths
	// This prevents high cardinality in metrics
	if len(path) > 0 && path[0] == '/' {
		parts := []rune(path)
		for i := 1; i < len(parts); i++ {
			if parts[i] == '/' {
				return string(parts[:i])
			}
		}
	}
	return path
}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	docs "github.com/xhkzeroone/go-generator/docs"
	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/handler"
	"github.com/xhkzeroone/go-generator/internal/middleware"
	"github.com/xhkzeroone/go-generator/internal/service"
)

// @title Go Generator API
// @version 1.0
// @description API endpoints for generating Go project scaffolding.
// @BasePath /
func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.InfoLevel)
	if os.Getenv("LOG_LEVEL") == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Initialize service
	manifestPath := constants.DefaultManifestPath
	if path := os.Getenv("MANIFEST_PATH"); path != "" {
		manifestPath = path
	}

	genService, err := service.NewGeneratorService(manifestPath)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize generator service")
	}

	// Initialize rate limiter (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute, logger)
	defer rateLimiter.Stop()

	// Initialize handlers
	genHandler := handler.NewGenerateHandler(genService, logger)
	healthHandler := handler.NewHealthHandler(logger)
	manifestHandler := handler.NewManifestHandler(genService, logger)

	// Setup routes
	mux := http.NewServeMux()

	docs.SwaggerInfo.BasePath = "/"

	// Apply middlewares (order matters: tracing -> metrics -> logging -> rate limit)
	tracingMiddleware := middleware.TracingMiddleware
	metricsMiddleware := middleware.MetricsMiddleware
	loggingMiddleware := middleware.LoggingMiddleware(logger)

	// Chain middlewares
	chainMiddleware := func(h http.Handler) http.Handler {
		return tracingMiddleware(metricsMiddleware(loggingMiddleware(h)))
	}

	// Metrics endpoint (no rate limiting, but with other middlewares)
	metricsHandler := handler.NewMetricsHandler()
	mux.Handle("/metrics", chainMiddleware(http.HandlerFunc(metricsHandler.HandleMetrics)))

	// API endpoints with all middlewares and rate limiting (must be before the catch-all)
	mux.Handle("/generate", chainMiddleware(rateLimiter.Limit(http.HandlerFunc(genHandler.HandleGenerate))))
	mux.Handle("/health", chainMiddleware(http.HandlerFunc(healthHandler.HandleHealth)))
	mux.Handle("/manifest", chainMiddleware(http.HandlerFunc(manifestHandler.HandleManifest)))
	mux.Handle("/swagger/", chainMiddleware(httpSwagger.WrapHandler))

	// Serve frontend as catch-all (with all middlewares)
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", chainMiddleware(fs))

	// Server configuration
	port := constants.DefaultPort
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.WithField("port", port).Info("Server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited")
}

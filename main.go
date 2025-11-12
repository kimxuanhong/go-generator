package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xhkzeroone/go-generator/internal/handler"
	"github.com/xhkzeroone/go-generator/internal/service"
)

func main() {
	// Initialize service
	manifestPath := "manifest.json"
	if path := os.Getenv("MANIFEST_PATH"); path != "" {
		manifestPath = path
	}

	genService, err := service.NewGeneratorService(manifestPath)
	if err != nil {
		log.Fatalf("Failed to initialize generator service: %v", err)
	}

	// Initialize handlers
	genHandler := handler.NewGenerateHandler(genService)
	healthHandler := handler.NewHealthHandler()
	manifestHandler := handler.NewManifestHandler(genService)

	// Setup routes
	mux := http.NewServeMux()

	// API endpoints (must be before the catch-all)
	mux.HandleFunc("/generate", genHandler.HandleGenerate)
	mux.HandleFunc("/health", healthHandler.HandleHealth)
	mux.HandleFunc("/manifest", manifestHandler.HandleManifest)

	// Serve frontend as catch-all
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	// Server configuration
	port := ":8080"
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
		log.Printf("Server starting on http://localhost%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

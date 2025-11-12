package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":  "ok",
		"service": "go-project-generator",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, we can't send JSON error, so just log it
		// In production, you might want to use a logger here
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *HealthHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

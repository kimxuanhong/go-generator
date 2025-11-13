package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse represents the payload returned by the health endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// HandleHealth godoc
// @Summary Check service health
// @Description Returns the health status of the generator service.
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 405 {object} ErrorResponse
// @Router /health [get]
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "go-project-generator",
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
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

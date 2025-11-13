package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xhkzeroone/go-generator/internal/models"
	"github.com/xhkzeroone/go-generator/internal/service"
)

type ManifestHandler struct {
	service *service.GeneratorService
}

func NewManifestHandler(svc *service.GeneratorService) *ManifestHandler {
	return &ManifestHandler{service: svc}
}

// HandleManifest godoc
// @Summary Get generator manifest
// @Description Returns the manifest describing available frameworks and libraries.
// @Tags manifest
// @Produce json
// @Success 200 {object} models.Manifest
// @Failure 405 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /manifest [get]
func (h *ManifestHandler) HandleManifest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	manifest := h.service.GetManifest()
	var manifestData *models.Manifest = manifest

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := json.NewEncoder(w).Encode(manifestData); err != nil {
		log.Printf("Error encoding manifest: %v", err)
		h.writeError(w, "Failed to encode manifest: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ManifestHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

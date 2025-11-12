package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xhkzeroone/go-generator/internal/service"
)

type ManifestHandler struct {
	service *service.GeneratorService
}

func NewManifestHandler(svc *service.GeneratorService) *ManifestHandler {
	return &ManifestHandler{service: svc}
}

func (h *ManifestHandler) HandleManifest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	manifest := h.service.GetManifest()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := json.NewEncoder(w).Encode(manifest); err != nil {
		log.Printf("Error encoding manifest: %v", err)
		h.writeError(w, "Failed to encode manifest: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ManifestHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

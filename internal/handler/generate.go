package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xhkzeroone/go-generator/internal/service"
)

type GenerateHandler struct {
	service *service.GeneratorService
}

func NewGenerateHandler(svc *service.GeneratorService) *GenerateHandler {
	return &GenerateHandler{service: svc}
}

// HandleGenerate godoc
// @Summary Generate a Go project scaffold
// @Description Generates a Go project scaffold based on the provided configuration and returns it as a ZIP archive.
// @Tags generator
// @Accept json
// @Produce application/zip
// @Param request body service.GenerateRequest true "Generator configuration"
// @Success 200 {file} file "Generated project archive"
// @Failure 400 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /generate [post]
func (h *GenerateHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req service.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.writeError(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	zipData, err := h.service.GenerateProject(&req)
	if err != nil {
		log.Printf("Error generating project: %v", err)
		h.writeError(w, "Failed to generate project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+req.ProjectName+".zip")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(zipData); err != nil {
		log.Printf("Error writing zip data: %v", err)
	}
}

func (h *GenerateHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

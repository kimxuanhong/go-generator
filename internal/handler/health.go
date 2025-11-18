package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/xhkzeroone/go-generator/internal/constants"
)

type HealthHandler struct {
	*BaseHandler
}

func NewHealthHandler(logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(logger),
	}
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
	if !h.validateMethod(r, constants.MethodGET) {
		h.writeError(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "go-project-generator",
	}

	h.writeJSON(w, http.StatusOK, response)
}

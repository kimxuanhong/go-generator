package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/xhkzeroone/go-generator/internal/constants"
	"github.com/xhkzeroone/go-generator/internal/service"
)

type ManifestHandler struct {
	*BaseHandler
	service *service.GeneratorService
}

func NewManifestHandler(svc *service.GeneratorService, logger *logrus.Logger) *ManifestHandler {
	return &ManifestHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     svc,
	}
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
	if !h.validateMethod(r, constants.MethodGET) {
		h.writeError(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	manifest := h.service.GetManifest()
	manifestData := manifest

	w.Header().Set(constants.HeaderCacheControl, constants.NoCache)
	w.Header().Set(constants.HeaderPragma, "no-cache")
	w.Header().Set(constants.HeaderExpires, "0")

	h.writeJSON(w, http.StatusOK, manifestData)
}

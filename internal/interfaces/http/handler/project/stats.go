package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *projectHandler) stats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	stats, err := h.projectService.GetProjectStats(ctx, projectID, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "fetch project stats failed", "error", err.Data)
		}
		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.SendResponse(w, meta.ReqID, http.StatusOK, "stats fetched", stats)
}

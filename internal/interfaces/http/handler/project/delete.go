package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *projectHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", "id")
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	err := h.projectService.DeleteProject(ctx, projectID, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "delete project", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Info(ctx, "delete project", "project_id", projectID)

	h.engine.SendNoContent(w)
}

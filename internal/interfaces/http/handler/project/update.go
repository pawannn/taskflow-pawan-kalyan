package projectHandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *projectHandler) update(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.engine.Log.Warn(ctx, apperr.ErrInvalidReqBody, "error", err)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrBadRequest, nil)
		return
	}

	updatedProject := models.Project{
		ID:          projectID,
		Name:        strings.TrimSpace(req.Name),
		Description: &req.Description,
	}

	project, err := h.projectService.UpdateProject(ctx, meta.UserID, updatedProject)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "update project", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Info(ctx, "update project", "project_id", projectID)

	h.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project updated", project)
}

package projectHandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pH *projectHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := pH.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		pH.engine.Log.Warn(ctx, domain.ErrValidationFailed, "fields", "id")
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pH.engine.Log.Warn(ctx, domain.ErrInvalidReqBody, "error", err)
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrBadRequest, nil)
		return
	}

	updatedProject := models.Project{
		ID:          projectID,
		Name:        strings.TrimSpace(req.Name),
		Description: &req.Description,
	}

	project, err := pH.projectService.UpdateProject(ctx, meta.UserID, updatedProject)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "update project", "error", err.Data)
		}

		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.Log.Info(ctx, "update project", "project_id", projectID)

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project updated", project)
}

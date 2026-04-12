package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pH *projectHandler) delete(w http.ResponseWriter, r *http.Request) {
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

	err := pH.projectService.DeleteProject(ctx, projectID, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "delete project", "error", err.Data)
		}

		pH.engine.SendResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.Log.Info(ctx, "delete project", "project_id", projectID)

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project deleted", nil)
}

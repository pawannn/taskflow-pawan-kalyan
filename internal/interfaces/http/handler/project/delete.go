package projectHandler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (pH *projectHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := pH.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if strings.TrimSpace(projectID) == "" {
		pH.engine.Log.Warn(ctx, "validation failed", "fields", "id")
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", map[string]string{
			"id": "is required",
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

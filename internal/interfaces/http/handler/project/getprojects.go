package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pH *projectHandler) projectByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(r.Context())

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		pH.engine.Log.Warn(ctx, domain.ErrValidationFailed, "fields", "id")
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	project, tasks, err := pH.projectService.GetProjectByID(ctx, projectID, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "fetch project", "error", err.Data)
		}

		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	resp := ProjectTasksResponse{
		Project: project,
		Tasks:   tasks,
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project fetched", resp)
}

func (pH *projectHandler) userProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(ctx)

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := utils.ParseIntDefault(pageStr, 1)
	limit := utils.ParseIntDefault(limitStr, 20)

	limit = min(20, limit)

	if page < 1 {
		page = 1
	}

	projects, err := pH.projectService.GetProjects(ctx, meta.UserID, page, limit)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "fetch user project", "error", err.Data)
		}

		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "projects fetched successfully", projects)
}

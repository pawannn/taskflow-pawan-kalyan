package projectHandler

import (
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
)

func (pH *projectHandler) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(r.Context())

	query := r.URL.Query()

	projectID := query.Get("id")
	if projectID == "" {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", map[string]string{
			"id": "is required",
		})
		return
	}

	project, tasks, err := pH.projectService.GetProjectByID(ctx, projectID)
	if !err.IsEmpty() {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusNotFound, domain.ErrNotFound, nil)
		return
	}

	resp := ProjectTasksResponse{
		Project: project,
		Tasks:   tasks,
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project fetched", resp)
}

func (pH *projectHandler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(ctx)

	if meta.UserID == "" {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusForbidden, domain.ErrForbidded, nil)
		return
	}

	projects, err := pH.projectService.GetProjects(ctx, meta.UserID)
	if !err.IsEmpty() {
		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "projects fetched successfully", projects)
}

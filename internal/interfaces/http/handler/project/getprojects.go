package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	taskHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/task"
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

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	limit, offset := utils.ParsePagination(pageStr, limitStr)

	project, tasks, hasNext, err := pH.projectService.GetProjectByID(ctx, projectID, meta.UserID, limit, offset)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "fetch project", "error", err.Data)
		}

		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	resp := ProjectTasksResponse{
		Project: project,
		Tasks: taskHandler.TasksResponse{
			Tasks: tasks,
			PaginationInfo: engine.PaginationInfo{
				Limit:   limit,
				Offset:  offset,
				HasNext: hasNext,
			},
		},
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project fetched", resp)
}

func (pH *projectHandler) userProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(ctx)

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	limit, offset := utils.ParsePagination(pageStr, limitStr)

	projects, hasNext, err := pH.projectService.GetProjects(ctx, meta.UserID, limit, offset)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "fetch user project", "error", err.Data)
		}

		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	response := userProjectsResponse{
		Projects: projects,
		PaginationInfo: engine.PaginationInfo{
			Limit:   limit,
			Offset:  offset,
			HasNext: hasNext,
		},
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "projects fetched successfully", response)
}

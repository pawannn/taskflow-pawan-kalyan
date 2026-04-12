package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	taskHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/task"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *projectHandler) projectByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := h.engine.ParseContext(r.Context())

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", "id")
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	limit, offset := utils.ParsePagination(pageStr, limitStr)

	project, tasks, hasNext, err := h.projectService.GetProjectByID(ctx, projectID, meta.UserID, limit, offset)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "fetch project", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
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

	h.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project fetched", resp)
}

func (h *projectHandler) userProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := h.engine.ParseContext(ctx)

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	limit, offset := utils.ParsePagination(pageStr, limitStr)

	projects, hasNext, err := h.projectService.GetProjects(ctx, meta.UserID, limit, offset)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "fetch user project", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
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

	h.engine.SendResponse(w, meta.ReqID, http.StatusOK, "projects fetched successfully", response)
}

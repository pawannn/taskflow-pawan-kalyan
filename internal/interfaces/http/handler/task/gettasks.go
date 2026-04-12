package taskHandler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tH *taskHandler) getByProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := tH.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		tH.engine.Log.Warn(ctx, domain.ErrValidationFailed, "fields", "id")
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")

	var taskStatus *models.TaskStatus
	var assigneePtr *string

	field := map[string]string{}

	if strings.TrimSpace(status) != "" {
		taskStats := models.TaskStatus(status)
		if !taskStats.IsValid() {
			field["status"] = "is invalid"
		} else {
			taskStatus = &taskStats
		}
	}

	if strings.TrimSpace(assignee) != "" {
		assigneePtr = &assignee
	}

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	limit, offset := utils.ParsePagination(pageStr, limitStr)

	tasks, hasNext, err := tH.taskService.GetByProjectID(ctx, projectID, taskStatus, assigneePtr, meta.UserID, limit, offset)
	if err != Error.NoErr {
		if err.Data != nil {
			tH.engine.Log.Error(ctx, "Fetch Task failed", "error", err.Data)
		}

		tH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	response := TasksResponse{
		Tasks: tasks,
		PaginationInfo: engine.PaginationInfo{
			Limit:   limit,
			Offset:  offset,
			HasNext: hasNext,
		},
	}

	tH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "tasks fetched", response)
}

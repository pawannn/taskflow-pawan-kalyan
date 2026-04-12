package taskHandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tH *taskHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := tH.engine.ParseContext(ctx)

	taskID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(taskID) {
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrBadRequest, nil)
		return
	}

	fields := make(map[string]string)
	var updatedTask models.Task

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			fields["title"] = "is invalid"
		} else {
			updatedTask.Title = title
		}
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if desc == "" {
			fields["description"] = "is invalid"
		} else {
			updatedTask.Description = &desc
		}
	}

	if req.Status != nil {
		s := models.TaskStatus(*req.Status)
		if !s.IsValid() {
			fields["status"] = "is invalid"
		} else {
			updatedTask.Status = s
		}
	}

	if req.Priority != nil {
		p := models.TaskPriority(*req.Priority)
		if !p.IsValid() {
			fields["priority"] = "is invalid"
		} else {
			updatedTask.Priority = &p
		}
	}

	if req.DueDate != nil {
		if strings.TrimSpace(*req.DueDate) == "" {
			fields["due_date"] = "is invalid"
		} else {
			d := utils.ParseDate(req.DueDate)
			updatedTask.DueDate = d
		}
	}

	if req.AssigneeID != nil {
		if !utils.IsValidUUID(*req.AssigneeID) {
			fields["assignee_id"] = "is invalid"
		} else {
			updatedTask.AssigneeID = req.AssigneeID
		}
	}

	if len(fields) > 0 {
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, fields)
		return
	}

	updatedTask.ID = taskID

	task, err := tH.taskService.UpdateTask(ctx, &updatedTask, meta.UserID)
	if err != Error.NoErr {
		tH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	tH.engine.Log.Info(ctx, "task updated", "task_id", taskID, "project_id", task.ProjectID)

	tH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "task updated", task)
}

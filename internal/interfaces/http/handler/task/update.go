package taskHandler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *taskHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	taskID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(taskID) {
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrBadRequest, nil)
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
		parsedDate := utils.ParseDate(req.DueDate)
		if parsedDate != nil && parsedDate.After(time.Now()) {
			updatedTask.DueDate = parsedDate
		} else {
			fields["due_date"] = "is invalid"
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
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, fields)
		return
	}

	updatedTask.ID = taskID

	task, err := h.taskService.UpdateTask(ctx, &updatedTask, meta.UserID)
	if err != apperr.NoErr {
		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Info(ctx, "task updated", "task_id", taskID, "project_id", task.ProjectID)

	h.engine.SendResponse(w, meta.ReqID, http.StatusOK, "task updated", task)
}

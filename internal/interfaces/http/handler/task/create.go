package taskHandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *taskHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := h.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", "id")
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.engine.Log.Warn(ctx, apperr.ErrInvalidReqBody, "error", err)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"title": req.Title,
	})

	var priority *models.TaskPriority
	if req.Priority != "" {
		p := models.TaskPriority(req.Priority)
		if !p.IsValid() {
			fields["priority"] = "is invalid"
		} else {
			priority = &p
		}
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		parsedDate := utils.ParseDate(req.DueDate)
		if parsedDate != nil && parsedDate.After(time.Now()) {
			dueDate = parsedDate
		} else {
			fields["due_date"] = "is invalid"
		}
	}

	if len(fields) > 0 {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", fields)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, fields)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: &req.Description,
		Priority:    priority,
		ProjectID:   projectID,
		DueDate:     dueDate,
	}

	err := h.taskService.CreateTask(ctx, &task, meta.UserID)

	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "create task failed", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Info(ctx, "task created", "task_id", task.ID, "project_id", projectID)

	h.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "task created", task)
}

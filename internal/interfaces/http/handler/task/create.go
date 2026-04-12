package taskHandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tH *taskHandler) create(w http.ResponseWriter, r *http.Request) {
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

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tH.engine.Log.Warn(ctx, domain.ErrInvalidReqBody, "error", err)
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
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

	if len(fields) > 0 {
		tH.engine.Log.Warn(ctx, domain.ErrValidationFailed, "fields", fields)
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, fields)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: &req.Description,
		Priority:    priority,
		ProjectID:   projectID,
		AssigneeID:  req.AssigneeID,
		DueDate:     utils.ParseDate(req.DueDate),
	}

	err := tH.taskService.CreateTask(ctx, &task, meta.UserID)

	if !err.IsEmpty() {
		if err.Data != nil {
			tH.engine.Log.Error(ctx, "create task failed", "error", err.Data)
		}

		tH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	tH.engine.Log.Info(ctx, "task created", "task_id", task.ID, "project_id", projectID)

	tH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "task created", task)
}

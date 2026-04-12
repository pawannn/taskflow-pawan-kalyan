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
		tH.engine.Log.Warn(ctx, "validation failed", "fields", "id")
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", map[string]string{
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

	priority := models.TaskPriority(req.Priority)
	if !priority.IsValid() {
		fields["priority"] = "is invalid"
	}

	if len(fields) > 0 {
		tH.engine.Log.Warn(ctx, "validation failed", "fields", fields)
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", fields)
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
			tH.engine.Log.Error(ctx, "create task", "error", err)
		}

		tH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	tH.engine.Log.Info(ctx, "task created", "task_id", task.ID, "project_id", projectID)

	tH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "task created", task)
}

package taskHandler

import (
	"encoding/json"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tH *taskHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := tH.engine.ParseContext(ctx)

	if meta.UserID == "" {
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusForbidden, domain.ErrForbidded, nil)
		return
	}

	query := r.URL.Query()
	projectID := query.Get("id")
	if projectID == "" {
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", map[string]string{
			"id": "is required",
		})
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	priority := models.TaskPriority(req.Priority)

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
		tH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	tH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "task created", task)
}

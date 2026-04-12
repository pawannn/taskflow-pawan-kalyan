package taskHandler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (tH *taskHandler) getByProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := tH.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")

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

	tasks, err := tH.taskService.GetTasks(ctx, projectID, taskStatus, assigneePtr, meta.UserID)
	if err != Error.NoErr {
		tH.engine.SendResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	tH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "tasks fetched", tasks)
}

package taskService

import (
	"context"
	"net/http"
	"time"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tS *TaskService) CreateTask(
	ctx context.Context,
	task *models.Task,
	userID string,
) Error.Err {
	if task.Title == "" {
		return Error.NewErr(http.StatusBadRequest, domain.ErrRequiredTaskTitle, nil)
	}

	if !task.Priority.IsValid() {
		return Error.NewErr(http.StatusBadRequest, domain.ErrInvalidTaskPriority, nil)
	}

	project, err := tS.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrFetchProject, nil)
	}

	if project == nil {
		return Error.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	if project.OwnerID != userID {
		return Error.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	task.ID = utils.GenerateUUID()
	task.Status = models.StatusTodo

	timeStamp := time.Now()
	task.CreatedAt = timeStamp
	task.UpdatedAt = timeStamp

	if err := tS.taskRepo.Create(ctx, task); err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrCreateTask, nil)
	}

	return Error.NoErr
}

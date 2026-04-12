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
	project, err := tS.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, nil)
	}

	if project == nil {
		return Error.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	if project.OwnerID != userID {
		return Error.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	task.ID = utils.GenerateUUID()
	task.Status = models.StatusTodo

	timeStamp := time.Now()
	task.CreatedAt = timeStamp
	task.UpdatedAt = timeStamp

	if err := tS.taskRepo.Create(ctx, task); err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, nil)
	}

	return Error.NoErr
}

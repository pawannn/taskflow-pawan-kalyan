package taskService

import (
	"context"
	"net/http"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tS *TaskService) CreateTask(
	ctx context.Context,
	task *models.Task,
	userID string,
) apperr.Err {
	project, err := tS.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, nil)
	}

	if project == nil {
		return apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	if project.OwnerID != userID {
		return apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
	}

	task.ID = utils.GenerateUUID()
	task.CreatorID = userID
	task.Status = models.StatusTodo

	timeStamp := time.Now()
	task.CreatedAt = timeStamp
	task.UpdatedAt = timeStamp

	if err := tS.taskRepo.Create(ctx, task); err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, nil)
	}

	return apperr.NoErr
}

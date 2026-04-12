package taskService

import (
	"context"
	"net/http"

	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (tS *TaskService) DeleteTask(ctx context.Context, taskID, userID string) apperr.Err {
	task, err := tS.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if task == nil {
		return apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	project, err := tS.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if project.OwnerID != userID && task.CreatorID != userID {
		return apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
	}

	err = tS.taskRepo.Delete(ctx, taskID)
	if err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return apperr.NoErr
}

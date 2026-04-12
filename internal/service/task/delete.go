package taskService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (tS *TaskService) DeleteTask(ctx context.Context, taskID, userID string) Error.Err {
	task, err := tS.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	if task == nil {
		return Error.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	project, err := tS.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	if project.OwnerID != userID {
		return Error.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	err = tS.taskRepo.Delete(ctx, taskID)
	if err != nil {
		return Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return Error.NoErr
}

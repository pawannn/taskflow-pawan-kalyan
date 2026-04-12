package taskService

import (
	"context"
	"net/http"
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (s *TaskService) UpdateTask(
	ctx context.Context,
	updatedTask *models.Task,
	userID string,
) (*models.Task, Error.Err) {
	task, err := s.taskRepo.GetByID(ctx, updatedTask.ID)
	if err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}
	if task == nil {
		return nil, Error.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	canUpdate, err := s.taskRepo.CanUpdateTask(ctx, updatedTask.ID, userID)
	if err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}
	if !canUpdate {
		return nil, Error.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	isUpdated := false

	if updatedTask.Title != "" && updatedTask.Title != task.Title {
		task.Title = updatedTask.Title
		isUpdated = true
	}

	if updatedTask.Description != nil {
		if task.Description == nil || *updatedTask.Description != *task.Description {
			task.Description = updatedTask.Description
			isUpdated = true
		}
	}

	if updatedTask.Status != "" && updatedTask.Status != task.Status {
		task.Status = updatedTask.Status
		isUpdated = true
	}

	if updatedTask.Priority != nil {
		if task.Priority == nil || *updatedTask.Priority != *task.Priority {
			task.Priority = updatedTask.Priority
			isUpdated = true
		}
	}

	if updatedTask.AssigneeID != nil {
		if task.AssigneeID == nil || *updatedTask.AssigneeID != *task.AssigneeID {
			user, err := s.userRepo.GetByID(ctx, *updatedTask.AssigneeID)
			if err != nil {
				return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
			}

			if user == nil {
				return nil, Error.NewErr(http.StatusBadRequest, domain.ErrBadRequest, err)
			}

			task.AssigneeID = updatedTask.AssigneeID
			isUpdated = true
		}
	}

	if updatedTask.DueDate != nil {
		if task.DueDate == nil || !updatedTask.DueDate.Equal(*task.DueDate) {
			task.DueDate = updatedTask.DueDate
			isUpdated = true
		}
	}

	if !isUpdated {
		return task, Error.NoErr
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return task, Error.NoErr
}

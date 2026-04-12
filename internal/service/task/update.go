package taskService

import (
	"context"
	"net/http"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (s *TaskService) UpdateTask(
	ctx context.Context,
	updatedTask *models.Task,
	userID string,
) (*models.Task, apperr.Err) {
	task, err := s.taskRepo.GetByID(ctx, updatedTask.ID)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}
	if task == nil {
		return nil, apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	canUpdate, err := s.taskRepo.CanUpdateTask(ctx, updatedTask.ID, userID)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}
	if !canUpdate {
		return nil, apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
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
				return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
			}

			if user == nil {
				return nil, apperr.NewErr(http.StatusBadRequest, apperr.ErrBadRequest, err)
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
		return task, apperr.NoErr
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return task, apperr.NoErr
}

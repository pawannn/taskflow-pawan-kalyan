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
	isTaskUpdated := false

	taskID := updatedTask.ID

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrFetchTask, err)
	}

	if task == nil {
		return nil, Error.NewErr(http.StatusNotFound, domain.ErrNotFound, err)
	}

	project, err := s.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrFetchProject, err)
	}

	if project.OwnerID != userID {
		return nil, Error.NewErr(http.StatusForbidden, domain.ErrForbidded, err)
	}

	if updatedTask.Title != task.Title {
		task.Title = updatedTask.Title
		isTaskUpdated = true
	}

	if updatedTask.Description != task.Description {
		task.Description = updatedTask.Description
		isTaskUpdated = true
	}

	if updatedTask.Status != task.Status {
		if !updatedTask.Status.IsValid() {
			return nil, Error.NewErr(http.StatusBadRequest, domain.ErrInvalidTaskStatus, err)
		}
		task.Status = updatedTask.Status
		isTaskUpdated = true
	}

	if updatedTask.Priority != task.Priority {
		if !updatedTask.Priority.IsValid() {
			return nil, Error.NewErr(http.StatusBadRequest, domain.ErrInvalidTaskPriority, err)
		}
		task.Priority = updatedTask.Priority
		isTaskUpdated = true
	}

	if updatedTask.AssigneeID != task.AssigneeID {
		task.AssigneeID = updatedTask.AssigneeID
		isTaskUpdated = true
	}

	if updatedTask.DueDate != task.DueDate {
		task.DueDate = updatedTask.DueDate
		isTaskUpdated = true
	}

	if !isTaskUpdated {
		return task, Error.NoErr
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrUpdateTask, err)
	}

	return task, Error.NoErr
}

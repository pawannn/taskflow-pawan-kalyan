package taskService

import (
	"context"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (s *TaskService) GetTasks(
	ctx context.Context,
	projectID string,
	status *models.TaskStatus,
	assigneeID *string,
	userID string,
) ([]*models.Task, Error.Err) {
	isauthorization, err := s.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if !isauthorization {
		return nil, Error.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	taskFilter := &domainRepo.TaskFilter{
		Status:     status,
		AssigneeID: assigneeID,
	}

	tasks, err := s.taskRepo.GetByProjectID(ctx, projectID, taskFilter, nil)
	if err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return tasks, Error.NoErr
}

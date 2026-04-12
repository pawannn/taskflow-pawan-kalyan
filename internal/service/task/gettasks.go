package taskService

import (
	"context"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (s *TaskService) GetByProjectID(
	ctx context.Context,
	projectID string,
	status *models.TaskStatus,
	assigneeID *string,
	userID string,
	limit, offset int,
) ([]*models.Task, bool, Error.Err) {
	pagination := domainRepo.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	isAuthorized, err := s.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if !isAuthorized {
		return nil, false, Error.NewErr(http.StatusForbidden, domain.ErrForbidden, nil)
	}

	taskFilter := &domainRepo.TaskFilter{
		Status:     status,
		AssigneeID: assigneeID,
	}

	tasks, hasNext, err := s.taskRepo.GetByProjectID(ctx, projectID, taskFilter, &pagination)
	if err != nil {
		return nil, false, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return tasks, hasNext, Error.NoErr
}

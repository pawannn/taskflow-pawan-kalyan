package taskService

import (
	"context"
	"net/http"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (s *TaskService) GetByProjectID(
	ctx context.Context,
	projectID string,
	status *models.TaskStatus,
	assigneeID *string,
	userID string,
	limit, offset int,
) ([]*models.Task, bool, apperr.Err) {
	pagination := domainRepo.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, false, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if project == nil {
		return nil, false, apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	isAuthorized, err := s.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if !isAuthorized {
		return nil, false, apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
	}

	taskFilter := &domainRepo.TaskFilter{
		Status:     status,
		AssigneeID: assigneeID,
	}

	tasks, hasNext, err := s.taskRepo.GetByProjectID(ctx, projectID, taskFilter, &pagination)
	if err != nil {
		return nil, false, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return tasks, hasNext, apperr.NoErr
}

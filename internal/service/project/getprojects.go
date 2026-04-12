package projectService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (pS *ProjectService) GetProjects(ctx context.Context, userID string, limit int, offset int) ([]*models.Project, bool, TaskFlowErr.Err) {
	pagination := domainRepo.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	projects, hasNext, err := pS.projectRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, false, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return projects, hasNext, TaskFlowErr.NoErr
}

func (pS *ProjectService) GetProjectByID(
	ctx context.Context,
	projectID string,
	userID string,
	limit, offset int,
) (*models.Project, []*models.Task, TaskFlowErr.Err) {
	pagination := domainRepo.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	isAuthorized, err := pS.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if !isAuthorized {
		return nil, nil, TaskFlowErr.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	project, err := pS.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	if project == nil {
		return nil, nil, TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	tasks, _, err := pS.taskRepo.GetByProjectID(ctx, projectID, nil, &pagination)
	if err != nil {
		return nil, nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return project, tasks, TaskFlowErr.NoErr
}

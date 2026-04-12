package projectService

import (
	"context"
	"net/http"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (pS *ProjectService) GetProjects(ctx context.Context, userID string, limit int, offset int) ([]*models.Project, bool, apperr.Err) {
	pagination := domainRepo.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	projects, hasNext, err := pS.projectRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, false, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return projects, hasNext, apperr.NoErr
}

func (pS *ProjectService) GetProjectByID(
	ctx context.Context,
	projectID string,
	userID string,
	limit, offset int,
) (*models.Project, []*models.Task, bool, apperr.Err) {
	pagination := domainRepo.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	project, err := pS.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, nil, false, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if project == nil {
		return nil, nil, false, apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	isAuthorized, err := pS.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if err != nil {
		return nil, nil, false, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if !isAuthorized {
		return nil, nil, false, apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
	}

	tasks, hasNext, err := pS.taskRepo.GetByProjectID(ctx, projectID, nil, &pagination)
	if err != nil {
		return nil, nil, false, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return project, tasks, hasNext, apperr.NoErr
}

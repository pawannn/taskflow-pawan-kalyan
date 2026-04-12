package projectService

import (
	"context"
	"net/http"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (pS *ProjectService) GetProjectStats(ctx context.Context, projectID, userID string) (*models.ProjectStats, apperr.Err) {
	isAuthorized, err := pS.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}
	if !isAuthorized {
		return nil, apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
	}

	project, err := pS.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}
	if project == nil {
		return nil, apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	stats, err := pS.taskRepo.GetProjectStats(ctx, projectID)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return stats, apperr.NoErr
}

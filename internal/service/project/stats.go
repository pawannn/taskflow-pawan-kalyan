package projectService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (pS *ProjectService) GetProjectStats(ctx context.Context, projectID, userID string) (*models.ProjectStats, TaskFlowErr.Err) {
	isAuthorized, err := pS.projectRepo.IsPartOfProject(ctx, projectID, userID)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}
	if !isAuthorized {
		return nil, TaskFlowErr.NewErr(http.StatusForbidden, domain.ErrForbidden, nil)
	}

	project, err := pS.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}
	if project == nil {
		return nil, TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	stats, err := pS.taskRepo.GetProjectStats(ctx, projectID)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return stats, TaskFlowErr.NoErr
}

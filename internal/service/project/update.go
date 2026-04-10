package projectService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (pS *ProjectService) UpdateProject(ctx context.Context, projectID, userID, name, description string) (*models.Project, TaskFlowErr.Err) {
	project, err := pS.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrFetchProject, err)
	}

	if project == nil {
		return nil, TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrProjectNotFound, nil)
	}

	if project.OwnerID != userID {
		return nil, TaskFlowErr.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	if name != "" {
		project.Name = name
	}

	if description != "" {
		project.Description = &description
	}

	if err := pS.projectRepo.Update(ctx, project); err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrUpdateProject, err)
	}

	return project, TaskFlowErr.NoErr
}

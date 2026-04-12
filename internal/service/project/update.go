package projectService

import (
	"context"
	"net/http"
	"time"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (pS *ProjectService) UpdateProject(ctx context.Context, ownerID string, updatedProject models.Project) (*models.Project, TaskFlowErr.Err) {
	needToUpdate := false

	project, err := pS.projectRepo.GetByID(ctx, updatedProject.ID)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	if project == nil {
		return nil, TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	if project.OwnerID != ownerID {
		return nil, TaskFlowErr.NewErr(http.StatusForbidden, domain.ErrForbidded, nil)
	}

	if updatedProject.Name != "" && project.Name != updatedProject.Name {
		project.Name = updatedProject.Name
		needToUpdate = true
	}

	if *updatedProject.Description != "" && *project.Description != *updatedProject.Description {
		project.Description = updatedProject.Description
		needToUpdate = true
	}

	if !needToUpdate {
		return project, TaskFlowErr.NoErr
	}

	project.UpdatedAt = time.Now()
	if err := pS.projectRepo.Update(ctx, project); err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return project, TaskFlowErr.NoErr
}

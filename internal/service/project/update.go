package projectService

import (
	"context"
	"net/http"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (pS *ProjectService) UpdateProject(ctx context.Context, ownerID string, updatedProject models.Project) (*models.Project, apperr.Err) {
	needToUpdate := false

	project, err := pS.projectRepo.GetByID(ctx, updatedProject.ID)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if project == nil {
		return nil, apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	if project.OwnerID != ownerID {
		return nil, apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
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
		return project, apperr.NoErr
	}

	project.UpdatedAt = time.Now()
	if err := pS.projectRepo.Update(ctx, project); err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return project, apperr.NoErr
}

package projectService

import (
	"context"
	"net/http"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pS *ProjectService) Create(ctx context.Context, name, description, ownerID string) (*models.Project, apperr.Err) {
	project := &models.Project{
		ID:          utils.GenerateUUID(),
		Name:        name,
		Description: &description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := pS.projectRepo.Create(ctx, project); err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return project, apperr.NoErr
}

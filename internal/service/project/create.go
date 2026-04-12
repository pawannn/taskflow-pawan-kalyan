package projectService

import (
	"context"
	"net/http"
	"time"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pS *ProjectService) Create(ctx context.Context, name, description, ownerID string) (*models.Project, Error.Err) {
	project := &models.Project{
		ID:          utils.GenerateUUID(),
		Name:        name,
		Description: &description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := pS.projectRepo.Create(ctx, project); err != nil {
		return nil, Error.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return project, Error.NoErr
}

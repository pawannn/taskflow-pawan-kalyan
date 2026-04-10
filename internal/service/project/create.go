package projectService

import (
	"context"
	"net/http"
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pS *ProjectService) Create(ctx context.Context, name, description, ownerID string) (*models.Project, TaskFlowErr.Err) {
	project := &models.Project{
		ID:          utils.GenerateUUID(),
		Name:        name,
		Description: &description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}

	if err := pS.projectRepo.Create(ctx, project); err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrCreateProject, err)
	}

	return project, TaskFlowErr.NoErr
}

package domainRepo

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

type Pagination struct {
	Limit  int
	Offset int
}

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetByUserID(ctx context.Context, userID string, pagination Pagination) ([]*models.Project, bool, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
	IsPartOfProject(ctx context.Context, projectID, userID string) (bool, error)
}

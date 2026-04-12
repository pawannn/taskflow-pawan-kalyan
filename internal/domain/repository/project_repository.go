// Package domainRepo contains repository contracts used by the domain layer.
package domainRepo

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

// Pagination represents limit and offset values for paginated queries.
type Pagination struct {
	Limit  int
	Offset int
}

// ProjectRepository defines data access methods for project persistence and retrieval.
type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetByUserID(ctx context.Context, userID string, pagination Pagination) ([]*models.Project, bool, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
	IsPartOfProject(ctx context.Context, projectID, userID string) (bool, error)
}

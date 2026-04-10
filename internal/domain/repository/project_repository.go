package repository

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
}

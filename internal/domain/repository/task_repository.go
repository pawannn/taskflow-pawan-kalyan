package domainRepo

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

type TaskFilter struct {
	Status     *models.TaskStatus
	AssigneeID *string
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	GetByProjectID(ctx context.Context, projectID string, filter *TaskFilter, pagination *Pagination) ([]*models.Task, error)
}

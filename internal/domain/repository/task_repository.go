package domainRepo

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

// TaskFilter represents filtering options for querying tasks.
type TaskFilter struct {
	Status     *models.TaskStatus
	AssigneeID *string
}

// TaskRepository defines data access methods for task persistence and queries.
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	CanUpdateTask(ctx context.Context, taskID string, userID string) (bool, error)
	GetByProjectID(ctx context.Context, projectID string, filter *TaskFilter, pagination *Pagination) ([]*models.Task, bool, error)
	GetProjectStats(ctx context.Context, projectID string) (*models.ProjectStats, error)
}

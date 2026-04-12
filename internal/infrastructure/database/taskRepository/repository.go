// Package taskRepository provides PostgreSQL implementations of task repository interfaces.
package taskRepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

// taskRepository implements TaskRepository using a PostgreSQL database.
type taskRepository struct {
	db *pgxpool.Pool
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(db *pgxpool.Pool) domainRepo.TaskRepository {
	return &taskRepository{db}
}

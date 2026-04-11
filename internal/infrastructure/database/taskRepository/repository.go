package taskRepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) domainRepo.TaskRepository {
	return &taskRepository{db}
}

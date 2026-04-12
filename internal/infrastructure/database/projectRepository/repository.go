// Package projectRepository provides PostgreSQL implementations of project repository interfaces.
package projectRepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

// projectRepository implements ProjectRepository using a PostgreSQL database.
type projectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new instance of ProjectRepository.
func NewProjectRepository(db *pgxpool.Pool) domainRepo.ProjectRepository {
	return &projectRepository{db}
}

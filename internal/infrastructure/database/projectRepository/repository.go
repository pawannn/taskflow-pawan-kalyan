package projectRepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

type projectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) domainRepo.ProjectRepository {
	return &projectRepository{db}
}

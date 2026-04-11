package projectRepository

import (
	"context"
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (pR *projectRepository) Create(ctx context.Context, project *models.Project) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `INSERT INTO projects (id, name, description, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := pR.db.Exec(ctx, query,
		project.ID,
		project.Name,
		project.Description,
		project.OwnerID,
		project.CreatedAt,
		project.UpdatedAt,
	)

	return err
}

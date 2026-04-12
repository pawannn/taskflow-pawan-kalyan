package projectRepository

import (
	"context"
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (r *projectRepository) Update(ctx context.Context, project *models.Project) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	project.UpdatedAt = time.Now()

	query := `UPDATE projects SET name = $1, description = $2, updated_at = $3 WHERE id = $4`

	_, err := r.db.Exec(ctx, query,
		project.Name,
		project.Description,
		project.UpdatedAt,
		project.ID,
	)

	return err
}

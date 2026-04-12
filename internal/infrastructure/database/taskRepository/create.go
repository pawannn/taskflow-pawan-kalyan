package taskRepository

import (
	"context"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

// Create inserts a new task record into the database.
func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO tasks (
			id, title, description, status, priority,
			project_id, assignee_id, creator_id, due_date,
			created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err := r.db.Exec(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.ProjectID,
		task.AssigneeID,
		task.CreatorID,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

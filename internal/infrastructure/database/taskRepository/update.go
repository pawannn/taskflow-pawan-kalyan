package taskRepository

import (
	"context"
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE tasks
		SET title=$1, description=$2, status=$3, priority=$4,
		    assignee_id=$5, due_date=$6, updated_at=$7
		WHERE id=$8
	`

	_, err := r.db.Exec(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.AssigneeID,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
	)

	return err
}

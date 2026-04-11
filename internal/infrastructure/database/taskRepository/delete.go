package taskRepository

import (
	"context"
	"time"
)

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

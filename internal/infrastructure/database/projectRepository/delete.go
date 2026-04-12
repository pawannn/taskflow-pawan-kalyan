package projectRepository

import (
	"context"
	"time"
)

// Delete removes a project record from the database by its ID.
func (r *projectRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `DELETE FROM projects WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

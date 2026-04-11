package projectRepository

import (
	"context"
	"time"
)

func (pR *projectRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `DELETE FROM projects WHERE id = $1`

	_, err := pR.db.Exec(ctx, query, id)
	return err
}

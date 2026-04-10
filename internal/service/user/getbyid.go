package userRepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

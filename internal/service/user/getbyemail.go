package userRepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

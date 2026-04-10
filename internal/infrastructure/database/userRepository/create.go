package userRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (uR *userRepository) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := uR.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("Unable to create user : %s", err.Error())
	}

	return nil
}

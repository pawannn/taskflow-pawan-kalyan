// Package userRepository provides PostgreSQL implementations of user repository interfaces.
package userRepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

// userRepository implements UserRepository using a PostgreSQL database.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *pgxpool.Pool) domainRepo.UserRepository {
	return &userRepository{db: db}
}

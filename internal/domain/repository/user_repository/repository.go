package userrepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/client"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) client.UserRepository {
	return &userRepository{db: db}
}

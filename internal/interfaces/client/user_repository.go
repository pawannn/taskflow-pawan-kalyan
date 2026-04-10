package client

import "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

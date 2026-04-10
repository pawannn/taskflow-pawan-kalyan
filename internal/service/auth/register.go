package authservice

import (
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	existing, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now()

	user := &models.User{
		ID:        utils.GenerateUUID(),
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, err
}

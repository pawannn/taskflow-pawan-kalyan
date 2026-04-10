package authservice

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// generate JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

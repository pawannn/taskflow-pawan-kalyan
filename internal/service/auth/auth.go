package authservice

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	Issuer    string
}

func NewAuthService(config *config.Config, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: config.JWTSecret,
		Issuer:    config.AppName,
	}
}

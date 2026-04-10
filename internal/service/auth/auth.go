package authservice

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastrcture/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/client"
)

type AuthService struct {
	userRepo  client.UserRepository
	jwtSecret string
	Issuer    string
}

func NewAuthService(config *config.Config, userRepo client.UserRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: config.JWTSecret,
		Issuer:    config.AppName,
	}
}

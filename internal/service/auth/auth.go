package service

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/client"
)

type AuthService struct {
	userRepo  client.UserRepository
	jwtSecret string
}

func NewAuthService(jwtSecret string, userRepo client.UserRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

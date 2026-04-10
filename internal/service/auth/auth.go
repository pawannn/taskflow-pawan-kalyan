package authService

import (
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
)

type AuthService struct {
	userRepo     domainRepo.UserRepository
	tokenService *auth.TokenService
	bCryptCost   int
}

func NewAuthService(
	bcryptCost int,
	userRepo domainRepo.UserRepository,
	tokenService *auth.TokenService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		bCryptCost:   bcryptCost,
		tokenService: tokenService,
	}
}

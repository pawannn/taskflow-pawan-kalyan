package authservice

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
)

type AuthService struct {
	userRepo     repository.UserRepository
	tokenService *auth.TokenService
	bCryptCost   int
}

func NewAuthService(bcryptCost int, userRepo repository.UserRepository, tokenService *auth.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		bCryptCost:   bcryptCost,
		tokenService: tokenService,
	}
}

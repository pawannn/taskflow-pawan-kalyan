package authHandler

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authservice "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
)

type AuthHandler struct {
	engine      *engine.HttpEngine
	authService *authservice.AuthService
}

func NewAuthHandler(engine *engine.HttpEngine, authService *authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		engine:      engine,
		authService: authService,
	}
}

func (e *AuthHandler) AddAuthRoutes() {
	e.engine.AddRoutes([]engine.Route{
		{
			Path:        "/auth/register",
			Method:      "POST",
			Description: "Register a user with name, email, password",
			Controller:  e.Register,
		},
		{
			Path:        "/auth/login",
			Method:      "POST",
			Description: "Authenticates a user using email and password, and returns a JWT access token on success.",
			Controller:  e.Login,
		},
	})
}

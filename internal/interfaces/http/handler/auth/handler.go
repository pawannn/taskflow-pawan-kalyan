package authhandler

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

func (e *AuthHandler) AddUserRoutes() {
	e.engine.AddRoutes([]engine.Route{
		{
			Path:        "/auth/register",
			Method:      "POST",
			Description: "Register a user with name, email, password",
			Controller:  e.Register,
		},
	})
}

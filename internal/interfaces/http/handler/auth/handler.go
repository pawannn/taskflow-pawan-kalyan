package authHandler

import (
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authservice "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
)

type authHandler struct {
	engine      *engine.HttpEngine
	authService *authservice.AuthService
}

func NewAuthHandler(engine *engine.HttpEngine, authService *authservice.AuthService) *authHandler {
	return &authHandler{
		engine:      engine,
		authService: authService,
	}
}

func (e *authHandler) AddRoutes() {
	e.engine.AddRoutes([]engine.Route{
		{
			Method:      http.MethodPost,
			Endpoint:    "/auth/register",
			Description: "Register a user with name, email, password",
			Controller:  e.Register,
		},
		{
			Method:      http.MethodPost,
			Endpoint:    "/auth/login",
			Description: "Returns a JWT access token",
			Controller:  e.Login,
		},
	})
}

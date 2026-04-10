package middlewares

import (
	"net/http"
	"strings"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authservice "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
)

type MiddlewareHandler struct {
	engine      *engine.HttpEngine
	authService *authservice.AuthService
}

func NewMiddlewareHadler(engine *engine.HttpEngine, authService *authservice.AuthService) *MiddlewareHandler {
	return &MiddlewareHandler{
		engine:      engine,
		authService: authService,
	}
}

func (m *MiddlewareHandler) ValidateAuthToken(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta := m.engine.ParseContext(r.Context())

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "authorization token not found", nil)
			return
		}

		Bearer, token, ok := strings.Cut(authHeader, " ")
		if !ok || Bearer != "Bearer" || token == "" {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "invalid authorization format", nil)
			return
		}

		claims, err := m.authService.ValidateJWT(token)
		if err != nil {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		meta.UserEmail = &claims.UserID
		meta.UserEmail = &claims.UserEmail

		m.engine.SetContext(r.Context(), meta)
	})
}

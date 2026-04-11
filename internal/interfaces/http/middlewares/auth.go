package middlewares

import (
	"net/http"
	"strings"

	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
)

type MiddlewareHandler struct {
	engine       *engine.HttpEngine
	tokenService auth.TokenService
}

func NewMiddlewareHadler(engine *engine.HttpEngine, tokenService auth.TokenService) *MiddlewareHandler {
	return &MiddlewareHandler{
		engine:       engine,
		tokenService: tokenService,
	}
}

func (m *MiddlewareHandler) ValidateAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta := m.engine.ParseContext(r.Context())

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "missing authorization token", nil)
			return
		}

		bearer, token, ok := strings.Cut(authHeader, " ")
		if !ok || bearer != "Bearer" || token == "" {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "invalid authorization format", nil)
			return
		}

		claims, err := m.tokenService.Validate(token)
		if err != nil {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		if claims.UserID == "" {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "invalid token: missing user_id", nil)
			return
		}

		if claims.UserEmail == "" {
			m.engine.SendResponse(w, meta.ReqID, http.StatusUnauthorized, "invalid token: missing user_email", nil)
			return
		}

		meta.UserID = claims.UserID
		meta.UserEmail = claims.UserEmail

		m.engine.SetContext(r.Context(), meta)

		next.ServeHTTP(w, r)
	})
}

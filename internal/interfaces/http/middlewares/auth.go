package middlewares

import (
	"net/http"
	"strings"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
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
		if strings.TrimSpace(authHeader) == "" {
			m.engine.SendErrorResponse(w, meta.ReqID, http.StatusUnauthorized, domain.ErrUnAuthorized, nil)
			return
		}

		bearer, token, ok := strings.Cut(authHeader, " ")
		if !ok || bearer != "Bearer" || token == "" {
			m.engine.SendErrorResponse(w, meta.ReqID, http.StatusUnauthorized, domain.ErrUnAuthorized, nil)
			return
		}

		claims, err := m.tokenService.Validate(token)
		if err != nil {
			m.engine.SendErrorResponse(w, meta.ReqID, http.StatusUnauthorized, domain.ErrUnAuthorized, nil)
			return
		}

		if strings.TrimSpace(claims.UserID) == "" {
			m.engine.SendErrorResponse(w, meta.ReqID, http.StatusUnauthorized, domain.ErrUnAuthorized, nil)
			return
		}

		if strings.TrimSpace(claims.UserEmail) == "" {
			m.engine.SendErrorResponse(w, meta.ReqID, http.StatusUnauthorized, domain.ErrUnAuthorized, nil)
			return
		}

		meta.UserID = claims.UserID
		meta.UserEmail = claims.UserEmail

		m.engine.SetContext(r.Context(), meta)

		next.ServeHTTP(w, r)
	})
}

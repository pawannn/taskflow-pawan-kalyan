package authhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (aH *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	aH.engine.Log.HTTP(ctx, r.Method, r.Pattern)
	meta := aH.engine.ParseContext(ctx)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		aH.engine.Log.Warn(ctx, "invalid request body", "error", err)
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"email":    req.Email,
		"password": req.Password,
	})

	if len(fields) > 0 {
		aH.engine.Log.Warn(ctx, "validation failed", "fields", fields)
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", map[string]interface{}{
			"fields": fields,
		})
		return
	}

	aH.engine.Log.Info(ctx, "login attempt", "email", req.Email)

	token, user, err := aH.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		errorMsg := "internal server error"
		statusCode := http.StatusInternalServerError

		if errors.Is(err, domain.ErrInvalidCredentials) {
			errorMsg = "invalid credentials"
			statusCode = http.StatusUnauthorized
			aH.engine.Log.Warn(ctx, "invalid credentials", "email", req.Email)
		} else {
			aH.engine.Log.Error(ctx, "login failed", "error", err)
		}

		aH.engine.SendResponse(w, meta.ReqID, statusCode, errorMsg, nil)
		return
	}

	aH.engine.Log.Info(ctx, "user logged in", "user_id", user.ID, "email", user.Email)

	response := AuthResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	aH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "login successful", response)
}

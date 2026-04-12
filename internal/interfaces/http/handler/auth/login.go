package authHandler

import (
	"encoding/json"
	"net/http"

	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.engine.Log.Warn(ctx, apperr.ErrInvalidReqBody, "error", err)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrInvalidReqBody, nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"email": req.Email,
	})

	if len(fields) > 0 {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", fields)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, fields)
		return
	}

	h.engine.Log.Auth(ctx, "login attempt", "email", req.Email)

	token, user, err := h.authService.Login(ctx, req.Email, req.Password)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "login failed", "email", req.Email, "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Auth(ctx, "user logged in", "user_id", user.ID, "email", user.Email)

	response := AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	h.engine.SendResponse(w, meta.ReqID, http.StatusOK, "login successful", response)
}

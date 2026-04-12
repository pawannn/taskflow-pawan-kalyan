package authHandler

import (
	"encoding/json"
	"net/http"

	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.engine.Log.Warn(ctx, "invalid request body", "error", err)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "invalid register details", nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"name":     req.Name,
		"email":    req.Email,
		"password": req.Password,
	})

	if len(fields) > 0 {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", fields)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, fields)
		return
	}

	user, err := h.authService.Register(ctx, req.Name, req.Email, req.Password)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "register failed", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Auth(ctx, "user registered", "user_id", user.ID, "email", user.Email)

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	h.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "user registered successfully", response)
}

package authhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (aH *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	aH.engine.Log.HTTP(ctx, r.Method, r.Pattern)
	meta := aH.engine.ParseContext(ctx)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		aH.engine.Log.Warn(ctx, "invalid request body", "error", err)
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid register details", nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"name":     req.Name,
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

	user, err := aH.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		errorMsg := "unable to register user"
		statusCode := http.StatusInternalServerError

		if errors.Is(err, domain.ErrUserAlreadyExists) {
			errorMsg = err.Error()
			statusCode = http.StatusConflict
		}

		aH.engine.Log.Error(ctx, "register failed", "error", err)
		aH.engine.SendResponse(w, meta.ReqID, statusCode, errorMsg, nil)
		return
	}

	aH.engine.Log.Info(ctx, "user registered", "user_id", user.ID, "email", user.Email)

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	aH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "user registered successfully", response)
}

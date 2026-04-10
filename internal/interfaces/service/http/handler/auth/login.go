package authhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/service/http/engine"
)

func (aH *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	meta := engine.ParseContext(r.Context())

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if req.Email == "" {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "email is required", nil)
		return
	}

	if req.Password == "" {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "password is required", nil)
		return
	}

	token, user, err := aH.authService.Login(req.Email, req.Password)
	if err != nil {
		errorMsg := "internal server error"
		statusCode := http.StatusInternalServerError

		if errors.Is(err, domain.ErrInvalidCredentials) {
			errorMsg = "invalid credentials"
			statusCode = http.StatusUnauthorized
		}

		aH.engine.SendResponse(w, meta.ReqID, statusCode, errorMsg, nil)
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	response := AuthResponse{
		Token: token,
		User:  userResponse,
	}

	aH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "login successful", response)
}

package authhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
)

func (aH *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	meta := aH.engine.ParseContext(r.Context())

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid register details", nil)
		return
	}

	if req.Name == "" {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid user name", nil)
		return
	}

	if req.Email == "" {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid user email", nil)
		return
	}

	if req.Password == "" {
		aH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid user password", nil)
		return
	}

	user, err := aH.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		errorMsg := "unable to register user"
		statusCode := http.StatusInternalServerError

		if errors.Is(err, domain.ErrUserAlreadyExists) {
			errorMsg = err.Error()
			statusCode = http.StatusConflict
		}

		aH.engine.SendResponse(w, meta.ReqID, statusCode, errorMsg, nil)
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	aH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "user registered successfully", response)
}

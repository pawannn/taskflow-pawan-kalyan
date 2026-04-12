package authService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	"golang.org/x/crypto/bcrypt"
)

func (aS *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, TaskFlowErr.Err) {
	user, err := aS.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	if user == nil {
		return "", nil, TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, TaskFlowErr.NewErr(http.StatusUnauthorized, domain.ErrUnAuthorized, nil)
	}

	token, err := aS.tokenService.Generate(user)
	if err != nil {
		return "", nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return token, user, TaskFlowErr.NoErr
}

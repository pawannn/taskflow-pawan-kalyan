package authService

import (
	"context"
	"net/http"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
)

func (aS *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, apperr.Err) {
	user, err := aS.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if user == nil {
		return "", nil, apperr.NewErr(http.StatusUnauthorized, apperr.ErrUnAuthorized, nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, apperr.NewErr(http.StatusUnauthorized, apperr.ErrUnAuthorized, nil)
	}

	token, err := aS.tokenService.Generate(user)
	if err != nil {
		return "", nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return token, user, apperr.NoErr
}

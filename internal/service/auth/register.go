package authService

import (
	"context"
	"net/http"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (aS *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, apperr.Err) {
	existing, err := aS.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, nil)
	}

	if existing != nil {
		return nil, apperr.NewErr(http.StatusConflict, apperr.ErrConflict, err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), int(aS.bCryptCost))
	if err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	timestamp := time.Now()

	user := &models.User{
		ID:        utils.GenerateUUID(),
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	if err := aS.userRepo.Create(ctx, user); err != nil {
		return nil, apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return user, apperr.NoErr
}

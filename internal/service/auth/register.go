package authService

import (
	"context"
	"net/http"
	"time"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (aS *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, TaskFlowErr.Err) {
	existing, err := aS.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, nil)
	}

	if existing != nil {
		return nil, TaskFlowErr.NewErr(http.StatusConflict, domain.ErrConflict, err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), int(aS.bCryptCost))
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
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
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return user, TaskFlowErr.NoErr
}

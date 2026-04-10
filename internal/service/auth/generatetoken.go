package authservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

type UserClaims struct {
	UserID    string
	UserEmail string
	jwt.RegisteredClaims
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := &UserClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "taskflow",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

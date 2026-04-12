package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

// Generate creates and signs a JWT token for the given user.
func (s *TokenService) Generate(user *models.User) (string, error) {
	claims := &UserClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.expiry) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

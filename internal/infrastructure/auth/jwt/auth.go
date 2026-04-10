package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID    string
	UserEmail string
	jwt.RegisteredClaims
}

type TokenService struct {
	issuer    string
	jwtSecret string
	expiry    int
}

func NewTokenService(issuer string, jwtSecret string, expiry int) *TokenService {
	return &TokenService{
		issuer:    issuer,
		jwtSecret: jwtSecret,
		expiry:    expiry,
	}
}

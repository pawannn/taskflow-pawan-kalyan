// Package auth provides utilities for authentication and JWT token handling.
package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// UserClaims represents custom JWT claims containing user identity information.
type UserClaims struct {
	UserID    string
	UserEmail string
	jwt.RegisteredClaims
}

// TokenService handles creation and management of JWT tokens.
type TokenService struct {
	issuer    string
	jwtSecret string
	expiry    int
}

// NewTokenService creates a new instance of TokenService.
func NewTokenService(issuer string, jwtSecret string, expiry int) *TokenService {
	return &TokenService{
		issuer:    issuer,
		jwtSecret: jwtSecret,
		expiry:    expiry,
	}
}

package domain

import (
	authdomain "portfolio-backend/internal/core/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenJwt     TokenType = "jwt"
	TokenRefresh TokenType = "refresh"
)

type TokenClaim struct {
	authdomain.User
	jwt.RegisteredClaims
}

package domain

import (
	authdomain "portfolio-backend/internal/core/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeJwt     TokenType = "jwt"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaim struct {
	authdomain.User
	jwt.RegisteredClaims
}

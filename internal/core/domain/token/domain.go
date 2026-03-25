package domain

import (
	authdomain "portfolio-backend/internal/core/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaim struct {
	authdomain.User
	jwt.RegisteredClaims
}

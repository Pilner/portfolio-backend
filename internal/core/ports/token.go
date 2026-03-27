package ports

import (
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
)

type TokenService interface {
	GenerateToken(tokenType tokendomain.TokenType, payload authdomain.User) (string, error)
	ValidateToken(tokenType tokendomain.TokenType, token string) (*authdomain.User, error)
}

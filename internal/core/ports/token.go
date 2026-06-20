package ports

import (
	authdomain "frv-backend/internal/core/domain/auth"
	tokendomain "frv-backend/internal/core/domain/token"
)

type TokenManager interface {
	GenerateToken(tokenType tokendomain.TokenType, payload authdomain.User) (string, error)
	ValidateToken(tokenType tokendomain.TokenType, token string) (*authdomain.User, error)
}

package service

import (
	authdomain "frv-backend/internal/core/domain/auth"
	tokendomain "frv-backend/internal/core/domain/token"
	"frv-backend/internal/core/ports"
)

type TokenService struct {
	tokenManager ports.TokenManager
}

func NewTokenManager(tokenManager ports.TokenManager) *TokenService {
	return &TokenService{tokenManager: tokenManager}
}

func (t *TokenService) GenerateToken(tokenType tokendomain.TokenType, payload authdomain.User) (string, error) {
	return t.tokenManager.GenerateToken(tokenType, payload)
}

func (t *TokenService) ValidateToken(tokenType tokendomain.TokenType, tokenString string) (*authdomain.User, error) {
	return t.tokenManager.ValidateToken(tokenType, tokenString)
}

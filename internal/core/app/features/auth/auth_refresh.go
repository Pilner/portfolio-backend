package features

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
)

type AuthRefreshHandler struct {
	tokenManager ports.TokenManager
}

func NewAuthRefreshHandler(tokenManager ports.TokenManager) AuthRefreshHandler {
	if tokenManager == nil {
		panic("nil token manager adapter")
	}
	return AuthRefreshHandler{
		tokenManager: tokenManager,
	}
}

func (h AuthRefreshHandler) Handle(ctx context.Context, userData *authdomain.User) (string, string, error) {
	accessToken, err := h.tokenManager.GenerateToken(tokendomain.TokenTypeJwt, *userData)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := h.tokenManager.GenerateToken(tokendomain.TokenTypeRefresh, *userData)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err

}

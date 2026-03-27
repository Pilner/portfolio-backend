package features

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
)

type AuthRefreshHandler struct {
	tokenService ports.TokenService
}

func NewAuthRefreshHandler(tokenService ports.TokenService) AuthRefreshHandler {
	if tokenService == nil {
		panic("nil token service")
	}
	return AuthRefreshHandler{
		tokenService: tokenService,
	}
}

func (h AuthRefreshHandler) Handle(ctx context.Context, userData *authdomain.User) (string, string, error) {
	accessToken, err := h.tokenService.GenerateToken(tokendomain.TokenTypeJwt, *userData)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := h.tokenService.GenerateToken(tokendomain.TokenTypeRefresh, *userData)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err

}

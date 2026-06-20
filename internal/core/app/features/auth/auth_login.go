package features

import (
	"context"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
)

type AuthLoginHandler struct {
	repo         authdomain.AuthRepository
	hashManager  ports.HashManager
	tokenManager ports.TokenManager
}

func NewAuthLoginHandler(repo authdomain.AuthRepository, hashManager ports.HashManager, tokenManager ports.TokenManager) AuthLoginHandler {
	if repo == nil {
		panic("nil auth repo")
	}
	if hashManager == nil {
		panic("nil hash manager adapter")
	}
	if tokenManager == nil {
		panic("nil token manager adapter")
	}
	return AuthLoginHandler{
		repo:         repo,
		hashManager:  hashManager,
		tokenManager: tokenManager,
	}
}

func (h AuthLoginHandler) Handle(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
	user, passwordHash, err := h.repo.FindUser(ctx, payload.Email)
	if err != nil {
		return user, "", "", err
	}

	if err := h.hashManager.Compare(passwordHash, payload.Password); err != nil {
		return user, "", "", domain.ErrPasswordDoesNotMatch
	}

	accessToken, err := h.tokenManager.GenerateToken(tokendomain.TokenTypeJwt, user)
	if err != nil {
		return user, "", "", err
	}

	refreshToken, err := h.tokenManager.GenerateToken(tokendomain.TokenTypeRefresh, user)
	if err != nil {
		return user, "", "", err
	}
	return user, accessToken, refreshToken, nil

}

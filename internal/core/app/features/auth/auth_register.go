package features

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
)

type AuthRegisterHandler struct {
	repo         authdomain.AuthRepository
	hashManager  ports.HashManager
	tokenManager ports.TokenManager
}

func NewAuthRegisterHandler(repo authdomain.AuthRepository, hashManager ports.HashManager, tokenManager ports.TokenManager) AuthRegisterHandler {
	if repo == nil {
		panic("nil auth repo")
	}
	if hashManager == nil {
		panic("nil hash manager adapter")
	}
	if tokenManager == nil {
		panic("nil token manager adapter")
	}
	return AuthRegisterHandler{
		repo:         repo,
		hashManager:  hashManager,
		tokenManager: tokenManager,
	}
}

func (h AuthRegisterHandler) Handle(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error) {
	hashedPassword, err := h.hashManager.Hash(payload.Password)
	if err != nil {
		return authdomain.User{}, "", "", err
	}

	payload.Password = hashedPassword

	user, err := h.repo.CreateUser(ctx, payload)
	if err != nil {
		return user, "", "", err
	}

	accessToken, err := h.tokenManager.GenerateToken(tokendomain.TokenTypeJwt, user)
	if err != nil {
		return user, "", "", err
	}

	refreshToken, err := h.tokenManager.GenerateToken(tokendomain.TokenTypeRefresh, user)
	if err != nil {
		return user, "", "", err
	}

	return user, accessToken, refreshToken, err
}

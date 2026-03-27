package features

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
)

type AuthRegisterHandler struct {
	repo         authdomain.AuthRepository
	hasher       ports.PasswordHasher
	tokenService ports.TokenService
}

func NewAuthRegisterHandler(repo authdomain.AuthRepository, hasher ports.PasswordHasher, tokenService ports.TokenService) AuthRegisterHandler {
	if repo == nil {
		panic("nil auth repo")
	}
	if hasher == nil {
		panic("nil password hasher")
	}
	if tokenService == nil {
		panic("nil token service")
	}
	return AuthRegisterHandler{
		repo:         repo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

func (h AuthRegisterHandler) Handle(ctx context.Context, payload authdomain.AddUser) (authdomain.User, string, string, error) {
	hashedPassword, err := h.hasher.Hash(payload.Password)
	if err != nil {
		return authdomain.User{}, "", "", err
	}

	payload.Password = hashedPassword

	user, err := h.repo.Register(ctx, payload)
	if err != nil {
		return user, "", "", err
	}

	accessToken, err := h.tokenService.GenerateToken(tokendomain.TokenJwt, user)
	if err != nil {
		return user, "", "", err
	}

	refreshToken, err := h.tokenService.GenerateToken(tokendomain.TokenRefresh, user)
	if err != nil {
		return user, "", "", err
	}

	return user, accessToken, refreshToken, err
}

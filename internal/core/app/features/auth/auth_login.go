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
	hasher       ports.PasswordHasher
	tokenService ports.TokenService
}

func NewAuthLoginHandler(repo authdomain.AuthRepository, hasher ports.PasswordHasher, tokenService ports.TokenService) AuthLoginHandler {
	if repo == nil {
		panic("nil auth repo")
	}
	if hasher == nil {
		panic("nil password hasher")
	}
	if tokenService == nil {
		panic("nil token service")
	}
	return AuthLoginHandler{
		repo:         repo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

func (h AuthLoginHandler) Handle(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
	user, passwordHash, err := h.repo.FindUser(ctx, payload.Email)
	if err != nil {
		return user, "", "", err
	}

	if err := h.hasher.Compare(passwordHash, payload.Password); err != nil {
		return user, "", "", domain.ErrPasswordDoesNotMatch
	}

	accessToken, err := h.tokenService.GenerateToken(tokendomain.TokenTypeJwt, user)
	if err != nil {
		return user, "", "", err
	}

	refreshToken, err := h.tokenService.GenerateToken(tokendomain.TokenTypeRefresh, user)
	if err != nil {
		return user, "", "", err
	}
	return user, accessToken, refreshToken, nil

}

package ports

import auth "portfolio-backend/internal/core/domain/auth"

type TokenService interface {
	GenerateAccessToken(payload auth.User) (string, error)
	GenerateRefreshToken(payload auth.User) (string, error)
	ValidateAccessToken(token string) (*auth.User, error)
	ValidateRefreshToken(token string) (*auth.User, error)
}

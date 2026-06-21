package service

import (
	"context"
	authdomain "frv-backend/internal/core/domain/auth"
	"frv-backend/internal/core/ports"
)

type AuthRepository struct {
	repository ports.AuthRepository
}

func NewAuthRepository(authRepository ports.AuthRepository) *AuthRepository {
	return &AuthRepository{repository: authRepository}
}

func (r *AuthRepository) CreateUser(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, error) {
	return r.repository.CreateUser(ctx, payload)
}

func (r *AuthRepository) FindUser(ctx context.Context, email string) (authdomain.User, string, error) {
	return r.repository.FindUser(ctx, email)
}

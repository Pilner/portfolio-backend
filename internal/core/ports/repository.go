package ports

import (
	"context"
	authdomain "frv-backend/internal/core/domain/auth"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, error)
	FindUser(ctx context.Context, email string) (authdomain.User, string, error)
}

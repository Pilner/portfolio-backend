package ports

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
)

type AuthRegister interface {
	Handle(ctx context.Context, payload authdomain.AddUser) (authdomain.User, string, string, error)
}

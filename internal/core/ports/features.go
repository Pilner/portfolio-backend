package ports

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
)

/* -------- V1 -------- */

// Auth
type AuthRegister interface {
	Handle(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error)
}
type AuthLogin interface {
	Handle(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error)
}

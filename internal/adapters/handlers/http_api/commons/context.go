package http_api

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
)

type AuthUserCtxKey struct{}

func AuthUserFromContext(ctx context.Context) (*authdomain.User, error) {
	u, ok := ctx.Value(AuthUserCtxKey{}).(*authdomain.User)
	if !ok {
		return u, ErrUserNotFoundInContext
	}

	return u, nil
}

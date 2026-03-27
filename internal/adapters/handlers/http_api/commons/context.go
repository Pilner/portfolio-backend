package http_api

import (
	"context"
	authdomain "portfolio-backend/internal/core/domain/auth"
)

type AuthUserCtxKey struct{}

func AuthUserFromContext(ctx context.Context) (*authdomain.User, bool) {
	u, ok := ctx.Value(AuthUserCtxKey{}).(*authdomain.User)
	return u, ok && u != nil
}

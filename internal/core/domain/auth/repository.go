package domain

import (
	"context"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, payload RegisterUser) (User, error)
	FindUser(ctx context.Context, email string) (User, string, error)
}

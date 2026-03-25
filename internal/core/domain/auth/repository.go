package domain

import (
	"context"
)

type AuthRepository interface {
	Register(ctx context.Context, payload AddUser) (User, error)
}

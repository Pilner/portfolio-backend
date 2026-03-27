package core

import (
	"context"
	"log/slog"
	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/adapters/crypto"
	"portfolio-backend/internal/adapters/repository"
	"portfolio-backend/internal/adapters/token"

	authfeature "portfolio-backend/internal/core/app/features/auth"
	"portfolio-backend/internal/core/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Features struct {
	AuthRegister ports.AuthRegister
	AuthLogin    ports.AuthLogin
}

type Application struct {
	Config   config.Values
	Features Features
	dbPool   *pgxpool.Pool
}

func NewApplication(ctx context.Context, config config.Values, logger *slog.Logger) Application {
	dbConfig, err := pgxpool.ParseConfig(config.DbConnectionUrl)
	if err != nil {
		panic(err)
	}
	dbConfig.MaxConns = 10
	postgresqlPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		panic(err)
	}
	err = postgresqlPool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	authRepo := repository.NewAuthPostgresRepository(postgresqlPool, logger)
	bcryptHasher := crypto.NewBcryptHasher()
	tokenService := token.NewJwtService(config)

	return Application{
		Features: Features{
			AuthRegister: authfeature.NewAuthRegisterHandler(authRepo, bcryptHasher, tokenService),
			AuthLogin:    authfeature.NewAuthLoginHandler(authRepo, bcryptHasher, tokenService),
		},
	}
}

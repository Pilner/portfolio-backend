package core

import (
	"context"
	"log/slog"
	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/adapters/crypto"
	authrepo "portfolio-backend/internal/adapters/repository/auth"
	"portfolio-backend/internal/adapters/token"
	"portfolio-backend/internal/service"

	authfeature "portfolio-backend/internal/core/app/features/auth"
	"portfolio-backend/internal/core/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Features struct {
	AuthRegister ports.AuthRegister
	AuthLogin    ports.AuthLogin
	AuthRefresh  ports.AuthRefresh
}

type Application struct {
	Config       config.Values
	Features     Features
	TokenManager ports.TokenManager
	dbPool       *pgxpool.Pool
}

func NewApplication(ctx context.Context, envConfig config.Values, logger *slog.Logger) Application {
	dbConfig, err := pgxpool.ParseConfig(envConfig.DbConnectionUrl)
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

	// Adapters
	authRepo := authrepo.NewAuthPostgresRepository(postgresqlPool, logger)
	bcryptHashManager := crypto.NewBcryptHashManager()
	jwtTokenManager := token.NewJwtTokenManager(envConfig)

	// Services
	hashManager := service.NewHashManager(bcryptHashManager)
	tokenManager := service.NewTokenManager(jwtTokenManager)

	return Application{
		Config: envConfig,
		Features: Features{
			AuthRegister: authfeature.NewAuthRegisterHandler(authRepo, hashManager, tokenManager),
			AuthLogin:    authfeature.NewAuthLoginHandler(authRepo, hashManager, tokenManager),
			AuthRefresh:  authfeature.NewAuthRefreshHandler(tokenManager),
		},
		TokenManager: tokenManager,
		dbPool:       postgresqlPool,
	}
}

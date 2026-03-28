package test

import (
	"context"
	"log/slog"
	"os"
	"portfolio-backend/internal/adapters/repository/migrations"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	once        sync.Once
	postgresUrl string
)

func GetPostgresUrl() string {
	return postgresUrl
}

func SetupDependencies(ctx context.Context) {
	once.Do(func() {
		SetupPostgres(ctx)
	})
}

func SetupPostgres(ctx context.Context) {
	postgresContainer, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase("frv"),
		postgres.WithPassword("secret"),
		postgres.WithUsername("root"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		panic(err)
	}
	postgresUrl, err = postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}
	migrator, err := migrations.NewMigrator(postgresUrl)
	if err != nil {
		panic(err)
	}
	err = migrator.Migrate()
	if err != nil {
		panic(err)
	}
}

func SetupTestLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

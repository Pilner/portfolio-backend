package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/adapters/handlers/http_api"
	"portfolio-backend/internal/adapters/repository/migrations"
	core "portfolio-backend/internal/core/app"
)

type CloseableService interface {
	Shutdown(ctx context.Context) error
}

func main() {
	// Load Timezone to calibrate server time
	loadTimezone()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	envConfig, err := config.LoadConfig(os.Getenv("ENV"), "./configs")
	if err != nil {
		logger.Error("configuration load error", "err", err)
		panic(err)
	}

	startMigration(envConfig, logger)

	application := core.NewApplication(context.TODO(), envConfig, logger)

	serviceCtx, serviceStopCtx := context.WithCancel(context.Background())

	httpServer := http_api.NewHttpApiServer(":8000", application, envConfig, logger)
	listenForShutdown(serviceCtx, serviceStopCtx, logger, httpServer)
	httpServer.StartServer()
}

func listenForShutdown(serviceCtx context.Context, serviceStopCtx context.CancelFunc, logger *slog.Logger, services ...CloseableService) {
	// Listen for syscall signals for process to interrupt/quit
	lifecycleLogger := logger.With("component", "AppLifecycle")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-sig
		lifecycleLogger.InfoContext(serviceCtx, "Got signal", "signal", s)
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownRelease := context.WithTimeout(serviceCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatalf("Graceful shutdown timed out... Forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		for _, service := range services {
			err := service.Shutdown(shutdownCtx)
			if err != nil {
				log.Fatal(err)
			}
		}

		shutdownRelease()
		serviceStopCtx()
	}()
}

func startMigration(config config.Values, logger *slog.Logger) {
	migratorLogger := logger.With("component", "DatabaseMigration")
	migrator, err := migrations.NewMigrator(config.DbConnectionUrl)
	if err != nil {
		panic(err)
	}
	defer func(migrator migrations.Migrator) {
		err := migrator.ReleaseConn()
		if err != nil {
			migratorLogger.Error("unable to release migration connection", "err", err)
		}
	}(migrator)

	// Get the current migration status
	now, exp, info, err := migrator.Info()
	if err != nil {
		panic(err)
	}

	if now < exp {
		migratorLogger.Info("migration needed, current state:", "info", info)

		err = migrator.Migrate()
		if err != nil {
			panic(err)
		}
		migratorLogger.Info("migration successful!")
	} else {
		migratorLogger.Info("no database migration needed")
	}
}

func loadTimezone() {
	// Set timezone
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		fmt.Println("Error loading time location", err)
	} else {
		time.Local = loc
	}

	currentTime := time.Now()
	fmt.Println("Loaded Timezone: ", time.Local.String())
	fmt.Println("Current Server Time: ", currentTime.Format(time.RFC3339))
}

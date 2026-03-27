package http_api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	core "portfolio-backend/internal/core/app"
	"time"

	"portfolio-backend/internal/adapters/handlers/http_api/middlewares"
	v1 "portfolio-backend/internal/adapters/handlers/http_api/v1"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	authRoutePath = "/auth"
)

type HttpApiServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewHttpApiServer(addr string, app core.Application, envCfg config.Values, logger *slog.Logger) HttpApiServer {

	rootRouter := chi.NewRouter()
	rootRouter.Use(middlewares.SetRequestId)

	adminRouter := v1.NewAdminRouter(app, logger)
	publicRouter := v1.NewPublicRouter(app, envCfg, logger)

	// Admin Routes
	rootRouter.Mount("/api/v1/admin", otelhttp.NewHandler(adminRouter, "admin-server"))

	// Public Routes
	rootRouter.Mount("/api/v1/", otelhttp.NewHandler(publicRouter, "public-server"))

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           rootRouter,
	}

	return HttpApiServer{
		server: server,
	}
}

func (s HttpApiServer) StartServer(port int) {
	fmt.Printf("Server starting on port %v\n", port)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s HttpApiServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	s.logger.InfoContext(ctx, "Stopped http api server")
	return nil
}

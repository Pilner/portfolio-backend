package http_api

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	core "portfolio-backend/internal/core/app"
	"time"

	auth "portfolio-backend/internal/adapters/handlers/http_api/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	authRoutePath = "/auth"
)

type HttpApiServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewHttpApiServer(addr string, app core.Application, envConfig config.Values, logger *slog.Logger) HttpApiServer {

	// Admin
	adminApiRouter := chi.NewRouter()
	adminApiRouter.Use(render.SetContentType(render.ContentTypeJSON))
	adminApiRouter.Use(otelchi.Middleware("frv-admin-service", otelchi.WithChiRoutes(adminApiRouter)))

	// Public
	authHandler := auth.NewPublicHandler(
		app.Features.AuthRegister,
		envConfig,
		logger,
	)

	publicApiRouter := chi.NewRouter()
	publicApiRouter.Use(render.SetContentType(render.ContentTypeJSON))
	publicApiRouter.Use(otelchi.Middleware("frv-public-service", otelchi.WithChiRoutes(publicApiRouter)))

	// Admin Routing

	// Public Routing
	publicApiRouter.Route(authRoutePath, func(r chi.Router) {
		r.Mount("/", authHandler.Routes())
	})

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/api/admin/v1", otelhttp.NewHandler(adminApiRouter, "admin-server"))
	rootRouter.Mount("/api/v1", otelhttp.NewHandler(publicApiRouter, "public-server"))

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           rootRouter,
	}

	return HttpApiServer{
		server: server,
	}
}

func (s HttpApiServer) StartServer() {
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

package v1

import (
	"log/slog"
	"portfolio-backend/internal/adapters/config"
	v1auth "portfolio-backend/internal/adapters/handlers/http_api/v1/auth"
	core "portfolio-backend/internal/core/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
)

func NewAdminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(otelchi.Middleware("frv-admin-service", otelchi.WithChiRoutes(r)))

	return r
}

func NewPublicRouter(app core.Application, envCfg config.Values, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(otelchi.Middleware("frv-public-service", otelchi.WithChiRoutes(r)))

	authHandler := v1auth.NewPublicHandler(
		app.Features.AuthRegister,
		app.Features.AuthLogin,
		envCfg,
		logger,
	)

	r.Route("/auth", func(r chi.Router) {
		r.Mount("/", authHandler.Routes())
	})

	return r
}

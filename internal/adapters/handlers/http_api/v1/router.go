package v1

import (
	"log/slog"
	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/adapters/handlers/http_api/middlewares"
	v1auth "portfolio-backend/internal/adapters/handlers/http_api/v1/auth"
	core "portfolio-backend/internal/core/app"
	tokendomain "portfolio-backend/internal/core/domain/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
)

func NewAdminRouter(app core.Application, envCfg config.Values, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(otelchi.Middleware("frv-admin-service", otelchi.WithChiRoutes(r)))

	// Auth
	authHandler := v1auth.NewAdminHandler(
		app.Features.AuthRefresh,
		envCfg,
		logger,
	)

	// Routes
	r.Route("/auth", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireToken(tokendomain.TokenTypeJwt, app.TokenService, logger))
			r.Get("/check", authHandler.AuthCheck)
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireToken(tokendomain.TokenTypeRefresh, app.TokenService, logger))
			r.Post("/refresh", authHandler.AuthRefresh)
		})

	})
	return r
}

func NewPublicRouter(app core.Application, envCfg config.Values, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(otelchi.Middleware("frv-public-service", otelchi.WithChiRoutes(r)))

	// Auth
	authHandler := v1auth.NewPublicHandler(
		app.Features.AuthRegister,
		app.Features.AuthLogin,
		envCfg,
		logger,
	)

	// Routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.AuthRegister)
		r.Post("/signin", authHandler.AuthLogin)
	})

	return r
}

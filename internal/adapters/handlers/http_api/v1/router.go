package v1

import (
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/adapters/handlers/http_api/middlewares"
	v1auth "portfolio-backend/internal/adapters/handlers/http_api/v1/auth"
	core "portfolio-backend/internal/core/app"
	tokendomain "portfolio-backend/internal/core/domain/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
)

func NewAdminRouter(app core.Application, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(otelchi.Middleware("frv-admin-service", otelchi.WithChiRoutes(r)))
	r.Use(middlewares.RequireToken(tokendomain.TokenTypeJwt, app.TokenService, logger))

	r.Get("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"message":"hello world"}`))
	})

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

package v1

import (
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	commons "portfolio-backend/internal/adapters/handlers/http_api/commons"
	"portfolio-backend/internal/adapters/handlers/http_api/middlewares"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	AuthCheckRenderErrFailed   = "auth check: rendering error failed"
	AuthRefreshRenderErrFailed = "auth refresh: rendering error failed"
)

type AdminHandler struct {
	authRefresh  ports.AuthRefresh
	tokenService ports.TokenService
	config       config.Values
	logger       *slog.Logger
}

func NewAdminHandler(
	authRefresh ports.AuthRefresh,
	tokenService ports.TokenService,
	cfg config.Values,
	logger *slog.Logger,
) AdminHandler {
	return AdminHandler{
		authRefresh:  authRefresh,
		tokenService: tokenService,
		config:       cfg,
		logger:       logger,
	}
}

func (h AdminHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.With(middlewares.RequireToken(tokendomain.TokenTypeJwt, h.tokenService, h.logger)).Get("/check", h.AuthCheck)
	r.With(middlewares.RequireToken(tokendomain.TokenTypeRefresh, h.tokenService, h.logger)).Post("/refresh", h.AuthRefresh)

	return r
}

func (h AdminHandler) AuthCheck(w http.ResponseWriter, r *http.Request) {
	user, ok := commons.AuthUserFromContext(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "auth check failed: no user found in context")
	}

	render.Status(r, http.StatusOK)
	render.Respond(w, r, UserResponse{
		Data: User{
			Id:          user.Id,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	})
}

func (h AdminHandler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	user, ok := commons.AuthUserFromContext(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "auth refresh failed: no user found in context")
	}

	accessToken, refreshToken, err := h.authRefresh.Handle(r.Context(), user)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "auth refresh failed", "error", err)
		apiError := commons.ErrInternalServerError(err)
		commons.RenderError(w, r, h.logger, AuthRefreshRenderErrFailed, apiError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.config.JwtTokenExpiryMinutes * 60, // MaxAge takes seconds
		HttpOnly: true,
		Secure:   h.config.IsProd, // Set to true in production
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   h.config.RefreshTokenExpiryMinutes * 60, // MaxAge takes seconds
		HttpOnly: true,
		Secure:   h.config.IsProd, // Set to true in production
		SameSite: http.SameSiteNoneMode,
	})

	h.logger.InfoContext(r.Context(), "access token refreshed successfully")

	render.Status(r, http.StatusOK)
	render.Respond(w, r, UserResponse{
		Data: User{
			Id:          user.Id,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	})
}

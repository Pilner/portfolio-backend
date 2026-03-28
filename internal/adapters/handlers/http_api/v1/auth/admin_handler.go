package v1

import (
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	commons "portfolio-backend/internal/adapters/handlers/http_api/commons"
	"portfolio-backend/internal/core/ports"

	"github.com/go-chi/render"
)

const (
	AuthCheckRenderErrFailed   = "auth check: rendering error failed"
	AuthRefreshRenderErrFailed = "auth refresh: rendering error failed"
)

type AdminHandler struct {
	authRefresh ports.AuthRefresh
	config      config.Values
	logger      *slog.Logger
}

func NewAdminHandler(
	authRefresh ports.AuthRefresh,
	envConfig config.Values,
	logger *slog.Logger,
) AdminHandler {
	return AdminHandler{
		authRefresh: authRefresh,
		config:      envConfig,
		logger:      logger,
	}
}

func (h AdminHandler) AuthCheck(w http.ResponseWriter, r *http.Request) {
	user, err := commons.AuthUserFromContext(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "auth check failed", "error", err)
		apiError := commons.ErrInternalServerError(err)
		commons.RenderError(w, r, h.logger, AuthCheckRenderErrFailed, apiError)
		return
	}

	h.logger.InfoContext(r.Context(), "auth checked", "user", user)

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
	user, err := commons.AuthUserFromContext(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "auth refresh failed", "error", err)
		apiError := commons.ErrInternalServerError(err)
		commons.RenderError(w, r, h.logger, AuthRefreshRenderErrFailed, apiError)
		return
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

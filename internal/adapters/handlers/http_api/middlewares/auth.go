package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	commons "portfolio-backend/internal/adapters/handlers/http_api/commons"
	"portfolio-backend/internal/core/domain"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
)

const (
	AuthLoginRenderErrFailed = "auth check token: rendering error failed"
)

// RequireToken validates the specified token cookie and, on success,
// stores the authenticated user in the request context.
func RequireToken(tokenType tokendomain.TokenType, tokenSvc ports.TokenService, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenCookie string
			switch tokenType {
			case tokendomain.TokenTypeJwt:
				tokenCookie = "access_token"
			case tokendomain.TokenTypeRefresh:
				tokenCookie = "refresh_token"
			default:
				logger.ErrorContext(r.Context(), "unknown token type", "tokenType", tokenType)
				commons.RenderError(w, r, logger, AuthLoginRenderErrFailed, commons.ErrUnauthorized())
				return
			}

			cookie, err := r.Cookie(tokenCookie)
			if err != nil || cookie.Value == "" {
				logger.ErrorContext(r.Context(), "missing token cookie", "cookie", tokenCookie, "error", err)
				commons.RenderError(w, r, logger, AuthLoginRenderErrFailed, commons.ErrUnauthorized())
				return
			}

			u, err := tokenSvc.ValidateToken(tokenType, cookie.Value)
			if err != nil {
				if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrExpiredToken) {
					logger.ErrorContext(r.Context(), "invalid token cookie", "cookie", tokenCookie, "error", err)
					commons.RenderError(w, r, logger, AuthLoginRenderErrFailed, commons.ErrUnauthorized())
					return
				}

				logger.ErrorContext(r.Context(), "token validation failed", "cookie", tokenCookie, "error", err)
				commons.RenderError(w, r, logger, AuthLoginRenderErrFailed, commons.ErrInternalServerError(err))
				return
			}

			// Set the user from the token to context
			ctx := context.WithValue(r.Context(), commons.AuthUserCtxKey{}, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

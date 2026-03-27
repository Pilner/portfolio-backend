package http_api

import (
	"errors"
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	errorsapi "portfolio-backend/internal/adapters/handlers/http_api/commons"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"
	"portfolio-backend/internal/core/ports"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	AuthRegisterRenderErrFailed = "auth register: rendering error failed"
)

type PublicHandler struct {
	authRegister ports.AuthRegister
	config       config.Values
	logger       *slog.Logger
}

func NewPublicHandler(
	authRegister ports.AuthRegister,
	cfg config.Values,
	logger *slog.Logger,
) PublicHandler {
	return PublicHandler{
		authRegister: authRegister,
		config:       cfg,
		logger:       logger,
	}
}

func (h PublicHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/signup", h.AuthRegister)

	return r
}

func (h PublicHandler) AuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		errorsapi.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, errorsapi.ErrEmptyRequest())
		return
	}

	req := RegisterAuth{}
	if err := render.Bind(r, &req); err != nil {
		err = errorsapi.TranslateBindError(err)
		apiError := errorsapi.ErrInvalidRequest(err, errorsapi.CodeInvalidRequest)
		errorsapi.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, apiError)
		return
	}

	payload := authdomain.AddUser{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	}

	data, accessToken, refreshToken, err := h.authRegister.Handle(r.Context(), payload)
	if err != nil {
		var apiError render.Renderer
		switch {
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			apiError = errorsapi.ErrConflict(err, errorsapi.CodeAuthEmailAlreadyExist)
		default:
			apiError = errorsapi.ErrInternalServerError(err)
		}

		h.logger.ErrorContext(r.Context(), "auth register user failed", "error", err)
		errorsapi.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, apiError)
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

	h.logger.InfoContext(r.Context(), "user registered successfully")

	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, User{
		Id:          data.Id,
		Email:       data.Email,
		DisplayName: data.DisplayName,
	})
}

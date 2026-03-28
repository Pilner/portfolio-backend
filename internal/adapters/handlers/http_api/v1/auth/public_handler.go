package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	commons "portfolio-backend/internal/adapters/handlers/http_api/commons"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"
	"portfolio-backend/internal/core/ports"

	"github.com/go-chi/render"
)

const (
	AuthRegisterRenderErrFailed = "auth register: rendering error failed"
	AuthLoginRenderErrFailed    = "auth login: rendering error failed"
)

type PublicHandler struct {
	authRegister ports.AuthRegister
	authLogin    ports.AuthLogin
	config       config.Values
	logger       *slog.Logger
}

func NewPublicHandler(
	authRegister ports.AuthRegister,
	authLogin ports.AuthLogin,
	envConfig config.Values,
	logger *slog.Logger,
) PublicHandler {
	return PublicHandler{
		authRegister: authRegister,
		authLogin:    authLogin,
		config:       envConfig,
		logger:       logger,
	}
}

func (h PublicHandler) AuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		commons.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, commons.ErrEmptyRequest())
		return
	}

	req := RegisterAuth{}
	if err := render.Bind(r, &req); err != nil {
		err = commons.TranslateBindError(err)
		apiError := commons.ErrInvalidRequest(err, commons.CodeInvalidRequest)
		commons.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, apiError)
		return
	}

	payload := authdomain.RegisterUser{
		AuthBase: authdomain.AuthBase{
			Email:    req.Email,
			Password: req.Password,
		},
		DisplayName: req.DisplayName,
	}

	data, accessToken, refreshToken, err := h.authRegister.Handle(r.Context(), payload)
	if err != nil {
		var apiError render.Renderer
		switch {
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			apiError = commons.ErrConflict(err, commons.CodeAuthEmailAlreadyExist)
		default:
			apiError = commons.ErrInternalServerError(err)
		}

		h.logger.ErrorContext(r.Context(), "auth register user failed", "error", err)
		commons.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, apiError)
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

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, UserResponse{
		Data: User{
			Id:          data.Id,
			Email:       data.Email,
			DisplayName: data.DisplayName,
		},
	})
}

func (h PublicHandler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		commons.RenderError(w, r, h.logger, AuthRegisterRenderErrFailed, commons.ErrEmptyRequest())
		return
	}

	req := LoginAuth{}
	if err := render.Bind(r, &req); err != nil {
		err = commons.TranslateBindError(err)
		apiError := commons.ErrInvalidRequest(err, commons.CodeInvalidRequest)
		commons.RenderError(w, r, h.logger, AuthLoginRenderErrFailed, apiError)
		return
	}

	payload := authdomain.LoginUser{
		AuthBase: authdomain.AuthBase{
			Email:    req.Email,
			Password: req.Password,
		},
	}

	data, accessToken, refreshToken, err := h.authLogin.Handle(r.Context(), payload)
	if err != nil {
		var apiError render.Renderer
		switch {
		case errors.Is(err, domain.ErrNoRecordsReturned),
			errors.Is(err, domain.ErrPasswordDoesNotMatch):
			apiError = commons.ErrUnauthorized()

		default:
			apiError = commons.ErrInternalServerError(err)
		}
		h.logger.ErrorContext(r.Context(), "auth login user failed", "error", err)
		commons.RenderError(w, r, h.logger, AuthLoginRenderErrFailed, apiError)
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

	h.logger.InfoContext(r.Context(), "user login successfully")

	render.Status(r, http.StatusOK)
	render.Respond(w, r, UserResponse{
		Data: User{
			Id:          data.Id,
			Email:       data.Email,
			DisplayName: data.DisplayName,
		},
	})
}

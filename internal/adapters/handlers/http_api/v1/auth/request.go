package v1

import (
	"net/http"
	commons "portfolio-backend/internal/adapters/handlers/http_api/commons"
)

type RegisterAuth struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (rb *RegisterAuth) Bind(r *http.Request) error {
	if rb.Email == "" {
		return commons.ErrEmptyEmail
	}
	normalizedEmail, err := commons.NormalizeEmailAlias(rb.Email)
	if err != nil {
		return commons.ErrInvalidEmail
	}
	rb.Email = normalizedEmail

	if rb.Password == "" {
		return commons.ErrEmptyPassword
	}
	if len(rb.Password) < 8 {
		return commons.ErrPasswordTooShort
	}

	if rb.DisplayName == "" {
		return commons.ErrEmptyDisplayName
	}

	if len(rb.DisplayName) > 50 {
		return commons.ErrInvalidDisplayName
	}

	return nil
}

type LoginAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rb *LoginAuth) Bind(r *http.Request) error {
	if rb.Email == "" {
		return commons.ErrEmptyEmail
	}

	if rb.Password == "" {
		return commons.ErrEmptyPassword
	}

	return nil
}

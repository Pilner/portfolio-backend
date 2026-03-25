package http_api

import (
	"net/http"
	errorsapi "portfolio-backend/internal/adapters/handlers/http_api/commons"
)

type RegisterAuth struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (rb *RegisterAuth) Bind(r *http.Request) error {
	if rb.Email == "" {
		return errorsapi.ErrEmptyUsername
	}

	if rb.Password == "" {
		return errorsapi.ErrEmptyPassword
	}

	if rb.DisplayName == "" {
		return errorsapi.ErrEmptyDisplayName
	}

	if len(rb.DisplayName) > 50 {
		return errorsapi.ErrInvalidDisplayName
	}

	return nil
}

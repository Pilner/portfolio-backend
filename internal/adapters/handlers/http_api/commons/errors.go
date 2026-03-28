package http_api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

/* API-Level Validation Errors */
var (
	ErrInvalidDataType       = errors.New("invalid data type in request payload")
	ErrEmptyEmail            = errors.New("email should have a value")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrEmptyPassword         = errors.New("password should have a value")
	ErrPasswordTooShort      = errors.New("password should be at least 8 characters long")
	ErrEmptyDisplayName      = errors.New("display_name should have a value")
	ErrInvalidDisplayName    = errors.New("invalid display_name format")
	ErrUserNotFoundInContext = errors.New("user not found in context")
)

/* Application Error Codes */
const (
	// General Errors (00)
	CodeInternalError  = "FRV000001"
	CodeNotFound       = "FRV000002"
	CodeEmptyBody      = "FRV000003"
	CodeInvalidRequest = "FRV000004"
	CodeUnauthorized   = "FRV000005"
	CodeForbidden      = "FRV000006"

	// Auth Errors (01)
	CodeAuthEmailAlreadyExist = "FRV010001"
)

type ApiError struct {
	Err            error `json:"-"`
	HttpStatusCode int   `json:"-"`

	Code      string         `json:"code"`
	Message   string         `json:"message"`
	RequestId string         `json:"requestId"`
	Timestamp string         `json:"timestamp"`
	Details   map[string]any `json:"details,omitempty"`
}

func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
	requestId := middleware.GetReqID(r.Context())
	e.RequestId = requestId
	e.Timestamp = time.Now().Local().Format(time.RFC3339)
	render.Status(r, e.HttpStatusCode)
	return nil
}

// HTTP 400
func ErrEmptyRequest() render.Renderer {
	return &ApiError{
		HttpStatusCode: http.StatusBadRequest,
		Code:           CodeEmptyBody,
		Message:        "Request Body Required",
	}
}

// HTTP 400
func ErrInvalidRequest(err error, code string) render.Renderer {
	return &ApiError{
		Err:            err,
		HttpStatusCode: http.StatusBadRequest,
		Code:           code,
		Message:        err.Error(),
	}
}

// HTTP 401
func ErrUnauthorized() render.Renderer {
	return &ApiError{
		HttpStatusCode: http.StatusUnauthorized,
		Code:           CodeUnauthorized,
		Message:        http.StatusText(http.StatusUnauthorized),
	}
}

// HTTP 403
func ErrForbidden(err error) render.Renderer {
	return &ApiError{
		Err:            err,
		HttpStatusCode: http.StatusForbidden,
		Code:           CodeForbidden,
		Message:        err.Error(),
	}
}

// HTTP 404
func ErrNotFound(err error) render.Renderer {
	return &ApiError{
		Err:            err,
		HttpStatusCode: http.StatusNotFound,
		Code:           CodeNotFound,
		Message:        err.Error(),
	}
}

// HTTP 409
func ErrConflict(err error, code string) render.Renderer {
	return &ApiError{
		Err:            err,
		HttpStatusCode: http.StatusConflict,
		Code:           code,
		Message:        err.Error(),
	}
}

// HTTP 500
func ErrInternalServerError(err error) render.Renderer {
	return &ApiError{
		Err:            err,
		HttpStatusCode: http.StatusInternalServerError,
		Code:           CodeInternalError,
		Message:        "Something went wrong. Please try again later.",
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, logMsg string, renderer render.Renderer) {
	if err := render.Render(w, r, renderer); err != nil {
		logger.ErrorContext(r.Context(), logMsg, "error", err)
	}
}

func TranslateBindError(err error) error {
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	if errors.As(err, &typeErr) || errors.As(err, &syntaxErr) {
		return ErrInvalidDataType
	}
	return err
}

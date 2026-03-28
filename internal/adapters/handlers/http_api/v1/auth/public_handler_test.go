package v1_test

import (
	"context"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	authhandler "portfolio-backend/internal/adapters/handlers/http_api/v1/auth"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"
	"portfolio-backend/internal/core/ports"
	"portfolio-backend/internal/test"
	"testing"
)

const publicApiAuthRoute = "/api/v1/auth"

type mockAuthRegisterFeature struct {
	mockHandle func(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error)
}

func (m *mockAuthRegisterFeature) Handle(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error) {
	return m.mockHandle(ctx, payload)
}

type mockAuthLoginFeature struct {
	mockHandle func(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error)
}

func (m *mockAuthLoginFeature) Handle(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
	return m.mockHandle(ctx, payload)
}

func TestHandlerAuthRegister(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description        string
		request            authhandler.RegisterAuth
		mockFeature        ports.AuthRegister
		expectedStatusCode int
	}{
		{
			description:        "failure - empty request body",
			request:            authhandler.RegisterAuth{},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - empty email",
			request: authhandler.RegisterAuth{
				Email:       "",
				Password:    "Test123456",
				DisplayName: "FooBar12",
			},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - invalid email",
			request: authhandler.RegisterAuth{
				Email:       "footest.com",
				Password:    "Test123456",
				DisplayName: "FooBar12",
			},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - empty password",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "",
				DisplayName: "FooBar12",
			},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - password too short",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "1234567",
				DisplayName: "FooBar12",
			},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - empty display name",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "Test123456",
				DisplayName: "",
			},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - display name too long",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "Test123456",
				DisplayName: "Lorem ipsum dolor sit amet, consectetuer adipiscina",
			},
			mockFeature:        &mockAuthRegisterFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - email already exists",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "Test123456",
				DisplayName: "FooBar12",
			},
			mockFeature: &mockAuthRegisterFeature{
				mockHandle: func(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error) {
					return authdomain.User{}, "", "", domain.ErrEmailAlreadyExists
				},
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			description: "failure - general error",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "Test123456",
				DisplayName: "FooBar12",
			},
			mockFeature: &mockAuthRegisterFeature{
				mockHandle: func(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error) {
					return authdomain.User{}, "", "", domain.ErrCreateUserFailed
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "success",
			request: authhandler.RegisterAuth{
				Email:       "foo@test.com",
				Password:    "Test123456",
				DisplayName: "FooBar12",
			},
			mockFeature: &mockAuthRegisterFeature{
				mockHandle: func(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, string, string, error) {
					return authdomain.User{
						Id:          "test-user-id",
						Email:       "foo@test.com",
						DisplayName: "FooBar12",
					}, "test-access-token", "test-refresh-token", nil
				},
			},
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := authhandler.NewPublicHandler(
				tc.mockFeature,
				&mockAuthLoginFeature{},
				config.Values{},
				test.SetupTestLogger(),
			)

			req := test.HttpRequest{
				T:                  t,
				Method:             http.MethodPost,
				Endpoint:           publicApiAuthRoute + "/signup",
				RequestBody:        tc.request,
				HandlerFunc:        handler.AuthRegister,
				ExpectedStatusCode: tc.expectedStatusCode,
			}

			resp := req.Execute()

			if tc.expectedStatusCode == http.StatusCreated {
				verifyCookie(t, resp.Cookies(), "test-access-token", "test-refresh-token")
			}
		})
	}
}

func TestHandlerAuthLogin(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description        string
		request            authhandler.LoginAuth
		mockFeature        ports.AuthLogin
		expectedStatusCode int
	}{
		{
			description:        "failure - empty request body",
			request:            authhandler.LoginAuth{},
			mockFeature:        &mockAuthLoginFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - empty email",
			request: authhandler.LoginAuth{
				Email:    "",
				Password: "Test123456",
			},
			mockFeature:        &mockAuthLoginFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - empty password",
			request: authhandler.LoginAuth{
				Email:    "foo@test.com",
				Password: "",
			},
			mockFeature:        &mockAuthLoginFeature{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "failure - no records found",
			request: authhandler.LoginAuth{
				Email:    "foo@test.com",
				Password: "Test123456",
			},
			mockFeature: &mockAuthLoginFeature{
				mockHandle: func(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
					return authdomain.User{}, "", "", domain.ErrNoRecordsReturned
				},
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			description: "failure - password mismatch",
			request: authhandler.LoginAuth{
				Email:    "foo@test.com",
				Password: "Test123456",
			},
			mockFeature: &mockAuthLoginFeature{
				mockHandle: func(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
					return authdomain.User{}, "", "", domain.ErrPasswordDoesNotMatch
				},
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			description: "failure - general error",
			request: authhandler.LoginAuth{
				Email:    "foo@test.com",
				Password: "Test123456",
			},
			mockFeature: &mockAuthLoginFeature{
				mockHandle: func(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
					return authdomain.User{}, "", "", domain.ErrFindUserFailed
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "success",
			request: authhandler.LoginAuth{
				Email:    "foo@test.com",
				Password: "Test123456",
			},
			mockFeature: &mockAuthLoginFeature{
				mockHandle: func(ctx context.Context, payload authdomain.LoginUser) (authdomain.User, string, string, error) {
					return authdomain.User{
						Id:          "test-user-id",
						Email:       "foo@test.com",
						DisplayName: "FooBar12",
					}, "test-access-token", "test-refresh-token", nil
				},
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := authhandler.NewPublicHandler(
				&mockAuthRegisterFeature{},
				tc.mockFeature,
				config.Values{},
				test.SetupTestLogger(),
			)

			req := test.HttpRequest{
				T:                  t,
				Method:             http.MethodPost,
				Endpoint:           publicApiAuthRoute + "/signin",
				RequestBody:        tc.request,
				HandlerFunc:        handler.AuthLogin,
				ExpectedStatusCode: tc.expectedStatusCode,
			}

			resp := req.Execute()

			if tc.expectedStatusCode == http.StatusOK {
				verifyCookie(t, resp.Cookies(), "test-access-token", "test-refresh-token")
			}
		})
	}
}

package v1_test

import (
	"context"
	"net/http"
	"portfolio-backend/internal/adapters/config"
	authhandler "portfolio-backend/internal/adapters/handlers/http_api/v1/auth"
	authdomain "portfolio-backend/internal/core/domain/auth"
	"portfolio-backend/internal/core/ports"
	"portfolio-backend/internal/test"
	"testing"
)

const adminApiAuthRoute = "/api/v1/admin/auth"

type mockAuthRefreshFeature struct {
	mockHandle func(ctx context.Context, userData *authdomain.User) (string, string, error)
}

func (m *mockAuthRefreshFeature) Handle(ctx context.Context, userData *authdomain.User) (string, string, error) {
	return m.mockHandle(ctx, userData)
}

func TestHandlerAuthCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description        string
		mockUserCtx        *authdomain.User
		expectedStatusCode int
	}{
		{
			description:        "failure - user not found in context",
			mockUserCtx:        nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "success",
			mockUserCtx: &authdomain.User{
				Id:          "c1e1c6b3-7c2f-4313-a4e8-a366103fe6be",
				Email:       "foo@test.com",
				DisplayName: "Test",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := authhandler.NewAdminHandler(
				&mockAuthRefreshFeature{},
				config.Values{},
				test.SetupTestLogger(),
			)

			req := test.HttpRequest{
				T:                  t,
				Method:             http.MethodGet,
				Endpoint:           adminApiAuthRoute + "/check",
				RequestBody:        nil,
				UserCtx:            tc.mockUserCtx,
				HandlerFunc:        handler.AuthCheck,
				ExpectedStatusCode: tc.expectedStatusCode,
			}

			req.Execute()
		})
	}
}

func TestHandlerAuthRefresh(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description          string
		mockFeature          ports.AuthRefresh
		mockUserCtx          *authdomain.User
		expectedStatusCode   int
		expectedJwtToken     string
		expectedRefreshToken string
	}{
		{
			description:          "failure - user not found in context",
			mockFeature:          &mockAuthRefreshFeature{},
			mockUserCtx:          nil,
			expectedJwtToken:     "",
			expectedRefreshToken: "",
			expectedStatusCode:   http.StatusInternalServerError,
		},
		{
			description: "success",
			mockFeature: &mockAuthRefreshFeature{
				mockHandle: func(ctx context.Context, userData *authdomain.User) (string, string, error) {
					return "test-access-token", "test-refresh-token", nil
				},
			},
			mockUserCtx: &authdomain.User{
				Id:          "c1e1c6b3-7c2f-4313-a4e8-a366103fe6be",
				Email:       "foo@test.com",
				DisplayName: "Test",
			},
			expectedJwtToken:     "test-access-token",
			expectedRefreshToken: "test-refresh-token",
			expectedStatusCode:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := authhandler.NewAdminHandler(
				tc.mockFeature,
				config.Values{
					JwtTokenExpiryMinutes:     30,
					RefreshTokenExpiryMinutes: 10080,
				},
				test.SetupTestLogger(),
			)

			req := test.HttpRequest{
				T:                  t,
				Method:             http.MethodPost,
				Endpoint:           adminApiAuthRoute + "/refresh",
				RequestBody:        nil,
				UserCtx:            tc.mockUserCtx,
				HandlerFunc:        handler.AuthRefresh,
				ExpectedStatusCode: tc.expectedStatusCode,
			}

			resp := req.Execute()

			if tc.mockUserCtx != nil {
				verifyCookie(t, resp.Cookies(), tc.expectedJwtToken, tc.expectedRefreshToken)
			}
		})
	}
}

// verifyCookie checks that both "access_token" and "refresh_token" cookies are present and have the expected values.
func verifyCookie(t *testing.T, cookies []*http.Cookie, expectedJwtToken, expectedRefreshToken string) {
	t.Helper()

	var foundAccessToken, foundRefreshToken bool
	for _, cookie := range cookies {
		switch cookie.Name {
		case "access_token":
			foundAccessToken = true
			if cookie.Value != expectedJwtToken {
				t.Errorf("access_token cookie "+test.TestUnexpectedValue.String(), expectedJwtToken, cookie.Value)
			}
		case "refresh_token":
			foundRefreshToken = true
			if cookie.Value != expectedRefreshToken {
				t.Errorf("refresh_token cookie "+test.TestUnexpectedValue.String(), expectedRefreshToken, cookie.Value)
			}
		}
	}

	if !foundAccessToken {
		t.Error("access_token cookie not set in response headers")
	}
	if !foundRefreshToken {
		t.Error("refresh_token cookie not set in response headers")
	}
}

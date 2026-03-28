package features_test

import (
	"context"
	"errors"
	features "portfolio-backend/internal/core/app/features/auth"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
	"portfolio-backend/internal/test"
	"testing"
)

type mockTokenService struct {
	mockGenerateToken func(tokenType tokendomain.TokenType, payload authdomain.User) (string, error)
	mockValidateToken func(tokenType tokendomain.TokenType, tokenString string) (*authdomain.User, error)
}

func (m *mockTokenService) GenerateToken(tokenType tokendomain.TokenType, payload authdomain.User) (string, error) {
	return m.mockGenerateToken(tokenType, payload)
}

func (m *mockTokenService) ValidateToken(tokenType tokendomain.TokenType, tokenString string) (*authdomain.User, error) {
	return m.mockValidateToken(tokenType, tokenString)
}

func testAuthRefresh(t *testing.T, tokenService ports.TokenService) {
	tokenGenerationError := errors.New("Fail")

	testCases := []struct {
		description        string
		payload            *authdomain.User
		shouldFailTokenGen bool
		expectedErr        error
	}{
		{
			description: "success",
			payload: &authdomain.User{
				Email:       "foo@test.com",
				DisplayName: "FooBar12",
			},
			expectedErr: nil,
		},
		{
			description: "failure - generate token failed",
			payload: &authdomain.User{
				Email:       "foo@test.com",
				DisplayName: "FooBar12",
			},
			shouldFailTokenGen: true,
			expectedErr:        tokenGenerationError,
		},
	}

	for _, tc := range testCases {
		var tokSvc ports.TokenService

		if tc.shouldFailTokenGen {
			tokSvc = &mockTokenService{
				mockGenerateToken: func(tokenType tokendomain.TokenType, payload authdomain.User) (string, error) {
					return "", tokenGenerationError
				},
			}
		} else {
			tokSvc = tokenService
		}

		feature := features.NewAuthRefreshHandler(tokSvc)
		t.Run(tc.description, func(t *testing.T) {
			_, _, err := feature.Handle(context.TODO(), tc.payload)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("error "+test.TestUnexpectedValue.String(), tc.expectedErr, err)
			}

		})
	}
}

package features_test

import (
	"context"
	"errors"
	features "portfolio-backend/internal/core/app/features/auth"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
	"portfolio-backend/internal/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testAuthRegister(t *testing.T, authRepo authdomain.AuthRepository, hasher ports.PasswordHasher, tokenService ports.TokenService) {
	tokenGenerationError := errors.New("Fail")

	testCases := []struct {
		description        string
		payload            authdomain.RegisterUser
		expectedData       authdomain.User
		shouldFailTokenGen bool
		expectedErr        error
	}{
		{
			description: "success",
			payload: authdomain.RegisterUser{
				AuthBase: authdomain.AuthBase{
					Email:    "foo@test.com",
					Password: "Test123456",
				},
				DisplayName: "FooBar12",
			},
			expectedData: authdomain.User{
				Email:       "foo@test.com",
				DisplayName: "FooBar12",
			},
			expectedErr: nil,
		},
		{
			description: "failure - email already exists",
			payload: authdomain.RegisterUser{
				AuthBase: authdomain.AuthBase{
					Email:    "foo@test.com",
					Password: "Test123456",
				},
				DisplayName: "FooBar12",
			},
			expectedData: authdomain.User{},
			expectedErr:  domain.ErrEmailAlreadyExists,
		},
		{
			description: "failure - generate token failed",
			payload: authdomain.RegisterUser{
				AuthBase: authdomain.AuthBase{
					Email:    "bar@test.com",
					Password: "Test123456",
				},
				DisplayName: "BarFoo12",
			},
			expectedData: authdomain.User{
				Email:       "bar@test.com",
				DisplayName: "BarFoo12",
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

		feature := features.NewAuthRegisterHandler(authRepo, hasher, tokSvc)

		t.Run(tc.description, func(t *testing.T) {
			user, _, _, err := feature.Handle(context.TODO(), tc.payload)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("error "+test.TestUnexpectedValue.String(), tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				// Don't include stochastic user id
				tc.expectedData.Id = user.Id

				if diff := cmp.Diff(tc.expectedData, user); diff != "" {
					t.Errorf("user data "+test.TestMismatchCompare.String(), diff)
				}
			}
		})
	}
}

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

func testAuthLogin(t *testing.T, authRepo authdomain.AuthRepository, hasher ports.PasswordHasher, tokenService ports.TokenService) {
	tokenGenerationError := errors.New("Fail")

	testCases := []struct {
		description        string
		payload            authdomain.LoginUser
		expectedData       authdomain.User
		shouldFailTokenGen bool
		expectedErr        error
	}{
		{
			description: "failure - no user found",
			payload: authdomain.LoginUser{
				AuthBase: authdomain.AuthBase{
					Email:    "INVALID",
					Password: "INVALID",
				},
			},
			expectedData: authdomain.User{},
			expectedErr:  domain.ErrNoRecordsReturned,
		},
		{
			description: "failure - password invalid",
			payload: authdomain.LoginUser{
				AuthBase: authdomain.AuthBase{
					Email:    "foo@test.com",
					Password: "INVALID",
				},
			},
			expectedData: authdomain.User{},
			expectedErr:  domain.ErrPasswordDoesNotMatch,
		},
		{
			description: "failure - generate token failed",
			payload: authdomain.LoginUser{
				AuthBase: authdomain.AuthBase{
					Email:    "foo@test.com",
					Password: "Test123456",
				},
			},
			shouldFailTokenGen: true,
			expectedData:       authdomain.User{},
			expectedErr:        tokenGenerationError,
		},
		{
			description: "success",
			payload: authdomain.LoginUser{
				AuthBase: authdomain.AuthBase{
					Email:    "foo@test.com",
					Password: "Test123456",
				},
			},
			expectedData: authdomain.User{
				Email:       "foo@test.com",
				DisplayName: "FooBar12",
			},
			expectedErr: nil,
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

		feature := features.NewAuthLoginHandler(authRepo, hasher, tokSvc)
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

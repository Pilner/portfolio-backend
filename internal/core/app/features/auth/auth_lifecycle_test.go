package features_test

import (
	"context"
	"log"
	"frv-backend/internal/adapters/config"
	"frv-backend/internal/adapters/crypto"
	authrepo "frv-backend/internal/adapters/repository/auth"
	"frv-backend/internal/adapters/token"
	"frv-backend/internal/service"
	"frv-backend/internal/test"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestAuthLifecycle(t *testing.T) {
	test.SetupPostgres(context.TODO())
	psqlPool, err := pgxpool.New(context.TODO(), test.GetPostgresUrl())
	if err != nil {
		log.Fatalf("fail to start new database pool: %v", err)
	}
	defer psqlPool.Close()

	envConfig := config.Values{
		JwtTokenExpiryMinutes:     30,
		RefreshTokenExpiryMinutes: 10080,
	}
	// Adapters
	authRepo := authrepo.NewAuthPostgresRepository(psqlPool, test.SetupTestLogger())
	bcryptHashManager := crypto.NewBcryptHashManager()
	jwtTokenManager := token.NewJwtTokenManager(envConfig)

	// Services
	hashManager := service.NewHashManager(bcryptHashManager)
	tokenManager := service.NewTokenManager(jwtTokenManager)

	t.Run("Test_AuthRegister", func(t *testing.T) {
		testAuthRegister(t, authRepo, hashManager, tokenManager)
	})
	t.Run("Test_AuthLogin", func(t *testing.T) {
		testAuthLogin(t, authRepo, hashManager, tokenManager)
	})
	t.Run("Test_AuthRefresh", func(t *testing.T) {
		testAuthRefresh(t, tokenManager)
	})
}

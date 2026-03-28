package features_test

import (
	"context"
	"log"
	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/adapters/crypto"
	"portfolio-backend/internal/adapters/repository"
	"portfolio-backend/internal/adapters/token"
	"portfolio-backend/internal/test"
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
	authRepo := repository.NewAuthPostgresRepository(psqlPool, test.SetupTestLogger())
	bcryptHasher := crypto.NewBcryptHasher()
	tokenService := token.NewJwtService(envConfig)

	t.Run("Test_AuthRegister", func(t *testing.T) {
		testAuthRegister(t, authRepo, bcryptHasher, tokenService)
	})
	t.Run("Test_AuthLogin", func(t *testing.T) {
		testAuthLogin(t, authRepo, bcryptHasher, tokenService)
	})
	t.Run("Test_AuthRefresh", func(t *testing.T) {
		testAuthRefresh(t, tokenService)
	})
}

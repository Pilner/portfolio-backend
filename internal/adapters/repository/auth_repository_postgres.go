package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable    = "users"
	userInfoTable = "user_info"
)

type AuthPostgresRepository struct {
	dbPool *pgxpool.Pool
	logger *slog.Logger
}

func NewAuthPostgresRepository(dbPool *pgxpool.Pool, logger *slog.Logger) AuthPostgresRepository {
	return AuthPostgresRepository{
		dbPool: dbPool,
		logger: logger.With("component", "AuthPostgresRepository"),
	}
}

func (r AuthPostgresRepository) Register(ctx context.Context, payload authdomain.AddUser) (authdomain.User, error) {
	u := authdomain.User{}

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to begin transaction for auth register", "error", err)
		return u, domain.ErrRegisterUserFailed
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			r.logger.ErrorContext(ctx, "failed to rollback transaction for auth register", "error", err)
		}
	}()

	userQuery := fmt.Sprintf(`
		INSERT INTO %v
		(
			id,
			email,
			password_hash
		)
		VALUES (gen_random_uuid(), $1, $2)
		RETURNING id, email
	`, usersTable)
	err = tx.QueryRow(
		ctx,
		userQuery,
		payload.Email,
		payload.Password,
	).Scan(&u.Id, &u.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
			return u, domain.ErrEmailAlreadyExists
		}

		r.logger.ErrorContext(ctx, "failed inserting to user table for auth register", "error", err)
		return u, domain.ErrRegisterUserFailed
	}

	userInfoQuery := fmt.Sprintf(`
		INSERT INTO %v
		(
			user_id,
			display_name
		)
		VALUES ($1, $2)
		RETURNING display_name
	`, userInfoTable)
	err = tx.QueryRow(ctx, userInfoQuery, u.Id, payload.DisplayName).Scan(&u.DisplayName)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed inserting to user_info table for auth register", "error", err)
		return u, domain.ErrRegisterUserFailed
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.ErrorContext(ctx, "failed to commit transaction for auth register", "error", err)
		return u, domain.ErrRegisterUserFailed
	}

	return u, nil
}

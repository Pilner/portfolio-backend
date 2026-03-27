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

func (r AuthPostgresRepository) CreateUser(ctx context.Context, payload authdomain.RegisterUser) (authdomain.User, error) {
	u := authdomain.User{}

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to begin transaction for auth register", "error", err)
		return u, domain.ErrCreateUserFailed
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
		return u, domain.ErrCreateUserFailed
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
		return u, domain.ErrCreateUserFailed
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.ErrorContext(ctx, "failed to commit transaction for auth register", "error", err)
		return u, domain.ErrCreateUserFailed
	}

	return u, nil
}

func (r AuthPostgresRepository) FindUser(ctx context.Context, email string) (authdomain.User, string, error) {
	u := authdomain.User{}
	var passwordHash string

	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.email,
			u.password_hash,
			ui.display_name
		FROM %v u
		INNER JOIN %v ui
			ON u.id = ui.user_id
		WHERE u.email = $1
	`, usersTable, userInfoTable)
	err := r.dbPool.QueryRow(ctx, query, email).Scan(
		&u.Id,
		&u.Email,
		&passwordHash,
		&u.DisplayName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, passwordHash, domain.ErrNoRecordsReturned
		}
		r.logger.ErrorContext(ctx, "failed fetching user for auth login", "error", err)
		return u, passwordHash, domain.ErrFindUserFailed
	}

	return u, passwordHash, nil
}

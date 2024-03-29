package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/igorakimy/grpc-sso-auth-service/internal/models"
	"github.com/igorakimy/grpc-sso-auth-service/internal/services/storage"
	"github.com/igorakimy/grpc-sso-auth-service/internal/storage/postgres/query"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(dsn string) (*Storage, func(), error) {
	const op = "storage.postgres.New"

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	closeConnFunc := func() {
		db.Close()
	}

	return &Storage{db: db}, closeConnFunc, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var userID int64
	err := s.db.QueryRow(ctx, query.InsertNewUser, email, passHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User

	err := s.db.QueryRow(ctx, query.SelectUserByEmail, email).
		Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool

	err := s.db.QueryRow(ctx, query.SelectIsAdminUserByID, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.IsAdmin"

	var app models.App

	err := s.db.QueryRow(ctx, query.SelectAppByID, appID).
		Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

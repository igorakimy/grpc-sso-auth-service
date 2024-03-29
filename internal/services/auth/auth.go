package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/igorakimy/grpc-sso-auth-service/internal/lib/jwt"
	"github.com/igorakimy/grpc-sso-auth-service/internal/services/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
)

type Service struct {
	log          *slog.Logger
	userProvider UserProvider
	userSaver    UserSaver
	appProvider  AppProvider
	ttl          time.Duration
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	userSaver UserSaver,
	appProvider AppProvider,
	ttl time.Duration,
) *Service {
	return &Service{
		log:          log,
		userProvider: userProvider,
		userSaver:    userSaver,
		appProvider:  appProvider,
		ttl:          ttl,
	}
}

func (s *Service) Login(ctx context.Context, email, password string, appID int32) (string, error) {
	const op = "auth.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("login user")

	user, err := s.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.log.Warn("user not found", slog.With(
				slog.String("op", op),
				slog.String("err", err.Error()),
			))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		s.log.Warn("failed to get user", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		log.Error("failed to compare password hashes", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := s.appProvider.App(ctx, int(appID))
	if err != nil {
		s.log.Error("failed to find application id", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewJWTToken(user, app, s.ttl)
	if err != nil {
		s.log.Error("failed to generate token", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully logged in")

	return token, nil
}
func (s *Service) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userId, err := s.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			s.log.Error("user already exists", slog.With(
				slog.String("op", op),
				slog.String("err", err.Error()),
			))
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		s.log.Error("failed to save user", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return userId, nil
}

func (s *Service) CheckIsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.CheckIsAdmin"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("check user is admin")

	isAdmin, err := s.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found", slog.With(
				slog.String("op", op),
				slog.String("err", err.Error()),
			))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		s.log.Info("failed to check is admin by user id", slog.With(
			slog.String("op", op),
			slog.String("err", err.Error()),
		))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

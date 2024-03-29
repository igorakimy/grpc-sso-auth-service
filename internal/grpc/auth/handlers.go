package auth

import (
	"context"
	"errors"
	"github.com/igorakimy/grpc-sso-auth-service/internal/services/auth"
	"github.com/igorakimy/grpc-sso-auth-service/internal/services/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/igorakimy/grpc-sso-auth-service/gen/protobuf/v1/sso"
)

func (s *server) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.authService.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "registration failed")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *server) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "invalid credentials")
	}

	return &ssov1.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *server) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.authService.CheckIsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidAppID) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to check")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "app ID is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == 0 {
		return status.Error(codes.Internal, "user id is required")
	}

	return nil
}

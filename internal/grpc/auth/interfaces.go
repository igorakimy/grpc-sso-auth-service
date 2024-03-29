package auth

import "context"

type AuthenticationService interface {
	Login(ctx context.Context, email, password string, appID int32) (string, error)
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	CheckIsAdmin(ctx context.Context, userID int64) (bool, error)
}

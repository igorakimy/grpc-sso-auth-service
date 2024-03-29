package suite

import (
	"context"
	"fmt"
	ssov1 "github.com/igorakimy/grpc-sso-auth-service/gen/protobuf/v1/sso"
	"github.com/igorakimy/grpc-sso-auth-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadFromPath("../config/test.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(
		ctx,
		fmt.Sprintf(":%d", cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

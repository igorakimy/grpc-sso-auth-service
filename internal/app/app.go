package app

import (
	"fmt"
	"github.com/igorakimy/grpc-sso-auth-service/internal/config"
	authGRPC "github.com/igorakimy/grpc-sso-auth-service/internal/grpc/auth"
	"github.com/igorakimy/grpc-sso-auth-service/internal/services/auth"
	"github.com/igorakimy/grpc-sso-auth-service/internal/storage/postgres"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	gRPCServer *grpc.Server
	cfg        *config.Config
	storage    *postgres.Storage
	log        *slog.Logger
}

func New(cfg *config.Config, storage *postgres.Storage, logger *slog.Logger) *App {
	gRPCServer := grpc.NewServer()

	authService := auth.New(logger, storage, storage, storage, cfg.TokenTTL)

	authGRPC.RegisterServer(gRPCServer, authService)

	return &App{
		gRPCServer: gRPCServer,
		cfg:        cfg,
		storage:    storage,
		log:        logger,
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	a.log.Info("application is running")

	log := a.log.With(
		slog.String("operation", op),
		slog.Int("port", a.cfg.Grpc.Port),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Grpc.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", lis.Addr().String()))

	if err = a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) GracefulShutdown() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	a.stopGrpcServer()

	a.log.Info("application stopped")
}

func (a *App) stopGrpcServer() {
	const operation = "app.Stop"

	a.log.With(slog.String("operation", operation)).
		Info("stopping gRPC server", slog.Int("port", a.cfg.Grpc.Port))

	a.gRPCServer.GracefulStop()
}

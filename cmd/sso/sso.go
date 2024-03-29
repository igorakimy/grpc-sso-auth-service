package main

import (
	"github.com/igorakimy/grpc-sso-auth-service/internal/app"
	"github.com/igorakimy/grpc-sso-auth-service/internal/config"
	"github.com/igorakimy/grpc-sso-auth-service/internal/lib/log"
	"github.com/igorakimy/grpc-sso-auth-service/internal/storage/postgres"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	logger := log.Setup(cfg.Env)

	// connect to storage
	storage, closeConn, err := postgres.New(cfg.Db.DSN)
	if err != nil {
		panic(err)
	}
	defer closeConn()

	// init app
	application := app.New(cfg, storage, logger)
	go application.MustRun()

	// graceful shutdown
	application.GracefulShutdown()
}

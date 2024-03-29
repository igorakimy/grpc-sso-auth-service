package main

import (
	"database/sql"
	"flag"
	"github.com/igorakimy/grpc-sso-auth-service/internal/config"
	"github.com/igorakimy/grpc-sso-auth-service/migrations"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var t string
	flag.StringVar(&t, "t", "up", "type of migration (up or down)")

	cfg := config.MustLoad()
	db, err := sql.Open("pgx", cfg.Db.DSN)
	if err != nil {
		panic(err)
	}
	defer func() { _ = db.Close() }()

	goose.SetBaseFS(migrations.EmbedMigrations)

	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if t == "up" {
		if err = goose.Up(db, "."); err != nil {
			panic(err)
		}
	} else if t == "down" {
		if err = goose.Down(db, "."); err != nil {
			panic(err)
		}
	}
}

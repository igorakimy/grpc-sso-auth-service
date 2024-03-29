-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(60) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE,
    secret VARCHAR(30) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE apps;
-- +goose StatementEnd

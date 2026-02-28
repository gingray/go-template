-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     email CITEXT NOT NULL UNIQUE,
                                     first_name TEXT NOT NULL,
                                     last_name TEXT NOT NULL,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

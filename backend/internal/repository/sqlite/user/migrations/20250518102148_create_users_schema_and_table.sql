-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS users;
CREATE TABLE IF NOT EXISTS users.users (
    id TEXT PRIMARY KEY,
    email_hash TEXT NOT NULL UNIQUE,
    email_encrypted BYTEA NOT NULL UNIQUE,
    password_hash TEXT,
    picture_encrypted BYTEA NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    logged_in_at TEXT NOT NULL,
    role TEXT NOT NULL,
    birthdate_encrypted BYTEA NOT NULL,
    lastname_encrypted BYTEA NOT NULL,
    firstname_encrypted BYTEA NOT NULL,
    gender_encrypted BYTEA NOT NULL,
    telephone_hash TEXT NOT NULL UNIQUE,
    telephone_encrypted BYTEA NOT NULL UNIQUE,
    postal_code_encrypted BYTEA NOT NULL,
    city_encrypted BYTEA NOT NULL,
    address1_encrypted BYTEA NOT NULL,
    address2_encrypted BYTEA,
    google_id_encrypted BYTEA,
    apple_id_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users.users;
DROP SCHEMA IF EXISTS users;
-- +goose StatementEnd

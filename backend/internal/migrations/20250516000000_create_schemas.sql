-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS votes;
CREATE SCHEMA IF NOT EXISTS settings;
CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS products;
CREATE SCHEMA IF NOT EXISTS messages;
CREATE SCHEMA IF NOT EXISTS events;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS votes;
DROP SCHEMA IF EXISTS settings;
DROP SCHEMA IF EXISTS users;
DROP SCHEMA IF EXISTS products;
DROP SCHEMA IF EXISTS messages;
DROP SCHEMA IF EXISTS events;
-- +goose StatementEnd

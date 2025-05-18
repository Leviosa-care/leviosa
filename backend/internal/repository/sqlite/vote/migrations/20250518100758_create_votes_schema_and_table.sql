-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS votes;

CREATE TABLE IF NOT EXISTS votes.votes (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    days TEXT NOT NULL CHECK (length(days) > 0),
    month INTEGER NOT NULL CHECK (month > 0 AND month < 13),
    year INTEGER NOT NULL CHECK (year >= 2024),

    UNIQUE (user_id, month, year),
    FOREIGN KEY (user_id) REFERENCES users.users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS votes.votes;
DROP SCHEMA IF EXISTS votes;
-- +goose StatementEnd

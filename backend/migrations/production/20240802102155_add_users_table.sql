-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email_hash TEXT NOT NULL UNIQUE,
    email_encrypted BLOB NOT NULL UNIQUE,
    password_hash TEXT,
    picture_encrypted BLOB,
	created_at TEXT NOT NULL,
	logged_in_at TEXT NOT NULL,
    role TEXT NOT NULL,
	birthdate_encrypted BLOB NOT NULL,
    lastname_encrypted BLOB NOT NULL,
    firstname_encrypted BLOB NOT NULL,
	gender_encrypted BLOB NOT NULL,
	telephone_hash TEXT NOT NULL UNIQUE,
	telephone_encrypted BLOB NOT NULL UNIQUE,
    postal_code_encrypted BLOB NOT NULL,
    city_encrypted BLOB NOT NULL,
    address1_encrypted BLOB NOT NULL,
    address2_encrypted BLOB,
    google_id_encrypted BLOB,
    apple_id_encrypted BLOB,
    dek_encrypted BLOB NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

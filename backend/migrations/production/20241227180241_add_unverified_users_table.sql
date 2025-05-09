-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS unverified_users (
    email_hash TEXT PRIMARY KEY,
    email_encrypted BLOB NOT NULL UNIQUE,
    password_hash BLOB NOT NULL UNIQUE,
	birthdate_encrypted  BLOB NOT NULL,
    lastname_encrypted BLOB NOT NULL,
    firstname_encrypted BLOB NOT NULL,
	gender_encrypted BLOB NOT NULL,
	telephone_hash TEXT NOT NULL UNIQUE,
	telephone_encrypted BLOB NOT NULL UNIQUE,
    created_at TEXT,
    postal_code_encrypted BLOB NOT NULL,
    city_encrypted BLOB NOT NULL,
    address1_encrypted BLOB NOT NULL,
    address2_encrypted BLOB,
    dek_encrypted BLOB NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE unverified_users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pending_users (
    id TEXT PRIMARY KEY,
    email_encrypted BLOB NOT NULL UNIQUE,
    picture_encrypted BLOB NOT NULL UNIQUE,
	birthdate_encrypted  BLOB NOT NULL,
    lastname_encrypted BLOB NOT NULL,
    firstname_encrypted BLOB NOT NULL,
	gender_encrypted BLOB NOT NULL,
	telephone_hash TEXT NOT NULL UNIQUE,
	telephone_encrypted BLOB NOT NULL UNIQUE,
    postal_code_encrypted BLOB NOT NULL UNIQUE,
    city_encrypted BLOB NOT NULL UNIQUE,
    address1_encrypted BLOB NOT NULL UNIQUE,
    address2_encrypted BLOB NOT NULL UNIQUE,
    google_id_encrypted BLOB,
    apple_id_encrypted BLOB,
    dek_encrypted BLOB NOT NULL UNIQUE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pending_users;
-- +goose StatementEnd

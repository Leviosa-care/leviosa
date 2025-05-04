-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    id,
    email_hash,
    email_encrypted,
    picture_encrypted,
    created_at,
    logged_in_at,
    birthdate_encrypted,
    lastname_encrypted,
    firstname_encrypted,
    gender_encrypted,
    telephone_hash,
    telephone_encrypted,
    postal_code_encrypted,
    city_encrypted,
    address1_encrypted,
    address2_encrypted,
    google_id_encrypted,
    apple_id_encrypted,
    dek_encrypted
) VALUES (
    "123e4567-e89b-12d3-a456-426614174000",
    'john.doe@example.com',
    'john.doe@example.com',
    'hashedpassword',
    'picture',
    '2025-02-03',
    '2025-02-03',
    'basic',
    '1998-07-12',
    'DOE',
    'John',
    'M',
    '0123456789',
    '0123456789',
    '75000',
    'Paris',
    '01 Avenue Jean DUPONT',
    '',
    'google_id',
    'apple_id',
    ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE rowid = 1;
-- +goose StatementEnd

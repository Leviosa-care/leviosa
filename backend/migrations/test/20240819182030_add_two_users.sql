-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    id,
    email_hash,
    email_encrypted,
    password_hash,
    picture_encrypted,
    created_at,
    logged_in_at,
    role,
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
) VALUES
    (
        "123e4567-e89b-12d3-a456-426614174001",
        'jane.doe@example.com',
        'jane.doe@example.com',
        'hashedpassword',
        'picture1',
        '2025-02-03',
        '2025-02-03',
        'basic',
        '1970-01-28',
        'DOE',
        'Jane',
        'F',
        '0123456781',
        '75000',
        'Paris',
        '01 Avenue Jean DUPONT',
        '',
        'google_id1',
        'apple_id1',
        ''
    ),
    (
        "123e4567-e89b-12d3-a456-426614174002",
        'jean.doe@example.com',
        'jean.doe@example.com',
        'hashedpassword',
        'picture2',
        '2025-02-03',
        '2025-02-03',
        'basic',
        '2000-10-05',
        'DOE',
        'Jean',
        'NB',
        '0123456782',
        '75000',
        'Paris',
        '01 Avenue Jean DUPONT',
        '',
        'google_id2',
        'apple_id2',
        ''
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE rowid > 2;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE settings.settings (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE settings.settings_encrypted (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    value_encrypted BYTEA NOT NULL,
    dek_encrypted BYTEA NOT NULL,
    key_version INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO settings.settings (key, value)
VALUES 
    ('company_name', 'leviosa'),
    ('support_email', 'contact@leviosa.care'),
    ('headquarters', '27 rue du Faubourg-Montmartre, 75009 Paris'),
    ('instagram_path', 'https://www.instagram.com/leviosa_care'),
    ('otp_duration', '15'),
    ('otp_length', '6'),
    ('otp_max_attempts', '3');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings.settings_encrypted;
DROP TABLE IF EXISTS settings.settings;
-- +goose StatementEnd

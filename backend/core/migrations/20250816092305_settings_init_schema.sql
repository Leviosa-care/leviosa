-- +goose Up
-- +goose StatementBegin

-- Create schema to group settings-related tables
CREATE SCHEMA IF NOT EXISTS settings;

-- ----------------------------------------------------------------------
-- Reusable trigger function to automatically update the 'updated_at' column
-- ----------------------------------------------------------------------
CREATE OR REPLACE FUNCTION settings.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ----------------------------------------------------------------------
-- Table: settings.plain
-- Stores plaintext configuration values
-- ----------------------------------------------------------------------
CREATE TABLE settings.plain (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,          -- Unique setting name
    value TEXT NOT NULL,               -- Plaintext value
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Attach trigger to auto-update the 'updated_at' timestamp
CREATE TRIGGER trg_update_config_timestamp
BEFORE UPDATE ON settings.plain
FOR EACH ROW
EXECUTE FUNCTION settings.update_timestamp();

-- ----------------------------------------------------------------------
-- Table: settings.encrypted
-- Stores encrypted configuration values
-- ----------------------------------------------------------------------
CREATE TABLE settings.encrypted (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,          -- Unique setting name
    value_encrypted BYTEA NOT NULL,    -- Encrypted value
    dek_encrypted BYTEA NOT NULL,      -- Data encryption key (encrypted)
    key_version INT NOT NULL,          -- Version of the encryption key used
    metadata JSONB,                    -- ENCX encryption metadata (KEK alias, timestamp, version)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Attach trigger to auto-update the 'updated_at' timestamp
CREATE TRIGGER trg_update_config_encrypted_timestamp
BEFORE UPDATE ON settings.encrypted
FOR EACH ROW
EXECUTE FUNCTION settings.update_timestamp();

-- ----------------------------------------------------------------------
-- Insert initial plaintext settings
-- ----------------------------------------------------------------------
INSERT INTO settings.plain (key, value) VALUES
    ('company_name', 'leviosa'),
    ('support_email', 'contact@leviosa.care'),
    ('headquarters', '27 rue du Faubourg-Montmartre, 75009 Paris'),
    ('instagram_path', 'https://www.instagram.com/leviosa_care'),
    ('otp_duration', '15'),   -- Duration in minutes
    ('otp_length', '6'),      -- Number of digits
    ('otp_max_attempts', '3'); -- Maximum allowed attempts

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings.encrypted;
DROP TABLE IF EXISTS settings.plain;
DROP FUNCTION IF EXISTS settings.update_timestamp();
DROP SCHEMA IF EXISTS settings;
-- +goose StatementEnd


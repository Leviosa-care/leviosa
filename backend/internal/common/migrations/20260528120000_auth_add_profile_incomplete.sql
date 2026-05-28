-- +goose Up
-- +goose StatementBegin

-- Add profile_incomplete column to track accounts created via guest booking claim.
-- These accounts have minimal data (name, email, phone, password) and need
-- profile completion (gender, birthdate, address) via an in-app nudge.
ALTER TABLE auth.users ADD COLUMN profile_incomplete BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE auth.users DROP COLUMN IF EXISTS profile_incomplete;

-- +goose StatementEnd

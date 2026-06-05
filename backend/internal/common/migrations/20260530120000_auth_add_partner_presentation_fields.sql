-- +goose Up
-- +goose StatementBegin

ALTER TABLE auth.partners
    ADD COLUMN occupation TEXT NOT NULL DEFAULT '',
    ADD COLUMN quote TEXT NOT NULL DEFAULT '',
    ADD COLUMN tags TEXT[] NOT NULL DEFAULT ARRAY[]::text[];

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE auth.partners
    DROP COLUMN IF EXISTS occupation,
    DROP COLUMN IF EXISTS quote,
    DROP COLUMN IF EXISTS tags;

-- +goose StatementEnd

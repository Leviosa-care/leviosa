-- +goose Up
-- +goose StatementBegin

ALTER TABLE auth.partners
    ADD COLUMN occupation TEXT,
    ADD COLUMN quote TEXT,
    ADD COLUMN tags TEXT[];

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE auth.partners
    DROP COLUMN IF EXISTS occupation,
    DROP COLUMN IF EXISTS quote,
    DROP COLUMN IF EXISTS tags;

-- +goose StatementEnd

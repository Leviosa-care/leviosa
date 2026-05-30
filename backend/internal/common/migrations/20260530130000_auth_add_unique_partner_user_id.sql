-- +goose Up
-- +goose StatementBegin

ALTER TABLE auth.partners
    ADD CONSTRAINT partners_user_id_unique UNIQUE (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE auth.partners
    DROP CONSTRAINT IF EXISTS partners_user_id_unique;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

ALTER TABLE auth.partners
    ADD COLUMN stripe_connected_account_id_encrypted BYTEA,
    ADD COLUMN stripe_account_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (stripe_account_status IN ('pending', 'active', 'restricted', 'disabled')),
    ADD COLUMN stripe_onboarding_complete BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_partners_stripe_account_status ON auth.partners (stripe_account_status);
CREATE INDEX idx_partners_stripe_onboarding_complete ON auth.partners (stripe_onboarding_complete);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS auth.idx_partners_stripe_account_status;
DROP INDEX IF EXISTS auth.idx_partners_stripe_onboarding_complete;

ALTER TABLE auth.partners
    DROP COLUMN IF EXISTS stripe_connected_account_id_encrypted,
    DROP COLUMN IF EXISTS stripe_account_status,
    DROP COLUMN IF EXISTS stripe_onboarding_complete;

-- +goose StatementEnd

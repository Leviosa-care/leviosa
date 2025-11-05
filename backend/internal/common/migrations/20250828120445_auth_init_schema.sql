-- +goose Up
-- +goose StatementBegin

-- Create auth schema
CREATE SCHEMA IF NOT EXISTS auth;

-- Create users table with encrypted/hashed fields only (GDPR compliant)
CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state VARCHAR(20) NOT NULL CHECK (state IN ('unverified', 'pending', 'active')),
    
    -- Email (encrypted + hashed for lookups)
    email_hash VARCHAR(64) NOT NULL UNIQUE,
    email_encrypted BYTEA NOT NULL,
    
    -- Password (hashed only for authentication)  
    password_hash_secure VARCHAR(255) NOT NULL,
    
    -- Role (encrypted for GDPR compliance)
    role_encrypted BYTEA,
    
    -- Profile data (all encrypted)
    picture_encrypted BYTEA,
    first_name_encrypted BYTEA,
    last_name_encrypted BYTEA,
    birth_date_encrypted BYTEA,
    gender_encrypted BYTEA,
    
    -- Contact info (telephone has both hash for lookup + encryption)
    telephone_hash VARCHAR(64),
    telephone_encrypted BYTEA,
    
    -- Address data (all encrypted)
    postal_code_encrypted BYTEA,
    city_encrypted BYTEA,
    address1_encrypted BYTEA,
    address2_encrypted BYTEA,
    
    -- External auth IDs (encrypted)
    google_id_encrypted BYTEA,
    apple_id_encrypted BYTEA,
    stripe_customer_id_encrypted BYTEA,
    
    -- Timestamps (encrypted for activity pattern protection)
    created_at_encrypted BYTEA,
    logged_in_at_encrypted BYTEA,
    
    -- Encryption metadata
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,
    
    -- Database metadata (unencrypted for system use)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for efficient lookups (only on hash fields, not encrypted)
CREATE INDEX idx_users_email_hash ON auth.users (email_hash);
CREATE INDEX idx_users_telephone_hash ON auth.users (telephone_hash) WHERE telephone_hash IS NOT NULL;
CREATE INDEX idx_users_state ON auth.users (state);
CREATE INDEX idx_users_created_at ON auth.users (created_at);
CREATE INDEX idx_users_google_id_encrypted ON auth.users (google_id_encrypted) WHERE google_id_encrypted IS NOT NULL;
CREATE INDEX idx_users_apple_id_encrypted ON auth.users (apple_id_encrypted) WHERE apple_id_encrypted IS NOT NULL;

-- Update trigger for updated_at timestamp
CREATE OR REPLACE FUNCTION auth.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON auth.users
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

-- Create partners table extending users with partner-specific data
CREATE TABLE auth.partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,

    -- Partner profile data (public information displayed to clients)
    bio TEXT,                                        -- Professional bio (public)
    experience TEXT,                                 -- Years of experience, background (public)
    -- certifications_encrypted BYTEA,                  -- Encrypted array of certifications - REMOVED

    -- Catalog associations (business metadata, not personal data)
    category_ids UUID[],                             -- Catalog category UUIDs this partner offers services for
    product_ids UUID[],                              -- Catalog product UUIDs this partner offers services for

    -- Stripe Connect integration (sensitive financial data remains encrypted)
    stripe_connected_account_id_encrypted BYTEA,     -- Encrypted Stripe Connected Account ID
    stripe_account_status VARCHAR(20) DEFAULT 'pending', -- Stripe account status (pending, active, restricted, disabled)
    stripe_onboarding_complete BOOLEAN DEFAULT FALSE, -- Whether Stripe onboarding is complete

    -- Encryption metadata (only needed for Stripe account ID)
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Database metadata (unencrypted for system use)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for partners table
CREATE INDEX idx_partners_id ON auth.partners (id);
CREATE INDEX idx_partners_user_id ON auth.partners (user_id);
CREATE INDEX idx_partners_stripe_account_status ON auth.partners (stripe_account_status);
CREATE INDEX idx_partners_stripe_onboarding_complete ON auth.partners (stripe_onboarding_complete);
CREATE INDEX idx_partners_created_at ON auth.partners (created_at);

CREATE TRIGGER update_partners_updated_at
    BEFORE UPDATE ON auth.partners
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

-- +goose StatementEnd

-- +goose Down  
-- +goose StatementBegin

-- Drop the auth schema and all its contents
DROP SCHEMA IF EXISTS auth CASCADE;

-- NOTE: should I keep that ?
DROP INDEX IF EXISTS auth.idx_users_google_id_encrypted;
DROP INDEX IF EXISTS auth.idx_users_apple_id_encrypted;
-- +goose StatementEnd

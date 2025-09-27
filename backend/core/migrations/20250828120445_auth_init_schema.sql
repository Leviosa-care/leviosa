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
    password_hash VARCHAR(255) NOT NULL,
    
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
    created_at_encrypted BYTEA NOT NULL,
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

-- Create specializations table for dynamic partner specialization types
CREATE TABLE auth.specializations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Specialization details (all encrypted for GDPR compliance)
    name_encrypted BYTEA NOT NULL,                    -- e.g., "physiotherapist", "mindset_coach"
    display_name_encrypted BYTEA NOT NULL,            -- e.g., "Physiotherapist", "Mindset Coach"
    description_encrypted BYTEA,                      -- Description of specialization

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Encryption metadata
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Database metadata (unencrypted for system use)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create partners table extending users with partner-specific data
CREATE TABLE auth.partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,

    -- Partner profile data (all encrypted for GDPR compliance)
    bio_encrypted BYTEA,                             -- Professional bio
    experience_encrypted BYTEA,                      -- Years of experience, background
    certifications_encrypted BYTEA,                  -- JSON array of certifications

    -- Verification status
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at_encrypted BYTEA,                     -- When verified (encrypted timestamp)
    verified_by_user_id UUID REFERENCES auth.users(id), -- Admin who verified

    -- Encryption metadata
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Database metadata (unencrypted for system use)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create junction table for partner specializations (many-to-many)
CREATE TABLE auth.partner_specializations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES auth.partners(id) ON DELETE CASCADE,
    specialization_id UUID NOT NULL REFERENCES auth.specializations(id) ON DELETE CASCADE,

    -- Database metadata (unencrypted for system use)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure unique combinations
    UNIQUE(partner_id, specialization_id)
);

-- Indexes for specializations table
CREATE INDEX idx_specializations_is_active ON auth.specializations (is_active);
CREATE INDEX idx_specializations_created_at ON auth.specializations (created_at);

-- Indexes for partners table
CREATE INDEX idx_partners_user_id ON auth.partners (user_id);
CREATE INDEX idx_partners_is_verified ON auth.partners (is_verified);
CREATE INDEX idx_partners_verified_by_user_id ON auth.partners (verified_by_user_id) WHERE verified_by_user_id IS NOT NULL;
CREATE INDEX idx_partners_created_at ON auth.partners (created_at);

-- Indexes for partner_specializations junction table
CREATE INDEX idx_partner_specializations_partner_id ON auth.partner_specializations (partner_id);
CREATE INDEX idx_partner_specializations_specialization_id ON auth.partner_specializations (specialization_id);
CREATE INDEX idx_partner_specializations_created_at ON auth.partner_specializations (created_at);

-- Update triggers for updated_at timestamps
CREATE TRIGGER update_specializations_updated_at
    BEFORE UPDATE ON auth.specializations
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

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

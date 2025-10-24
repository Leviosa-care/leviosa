-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TYPE catalog.status AS ENUM ('draft', 'published', 'archived');
CREATE TYPE catalog.price_interval AS ENUM ('month', 'year', 'one_time');
CREATE TYPE catalog.coupon_duration AS ENUM ('once', 'repeating', 'forever');

CREATE TABLE IF NOT EXISTS catalog.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- <-- Change from TEXT to UUID
    name TEXT NOT NULL UNIQUE, -- Maps to 'Name string'. Name should be unique.
    description TEXT, -- Maps to 'Description string,omitempty'. 'omitempty' implies NULLable. Removed UNIQUE as descriptions are not typically unique.
    status catalog.status NOT NULL DEFAULT 'draft', -- New field with ENUM type and default
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB -- Maps to 'Metadata map[string]any,omitempty'. JSONB for flexible key-value pairs. 'omitempty' implies NULLable.
);

-- Trigger to update updated_at automatically (useful for audit/caching)
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_categories_updated_at
BEFORE UPDATE ON catalog.categories
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Add comments for better documentation (optional but good practice)
COMMENT ON TABLE catalog.categories IS 'Defines categories for products, allowing for flexible grouping and custom attributes.';
COMMENT ON COLUMN catalog.categories.id IS 'Unique identifier for the category (e.g., UUID).';
COMMENT ON COLUMN catalog.categories.name IS 'Name of the category (must be unique).';
COMMENT ON COLUMN catalog.categories.description IS 'Description of the category.';
COMMENT ON COLUMN catalog.categories.metadata IS 'Flexible JSONB field for additional, category-specific attributes or configurations.';

CREATE TABLE IF NOT EXISTS catalog.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE, -- Maps to 'Name string'. Keeping UNIQUE as product names are often unique.
    description TEXT NOT NULL, -- Maps to 'Description string,omitempty'. 'omitempty' means it can be NULL. Removed UNIQUE as descriptions are rarely unique.
    category_id UUID NOT NULL, -- Maps to 'CategoryID string'. Using snake_case. Removed UNIQUE, as many products can share a category.
    duration INT NOT NULL, -- Maps to 'Duration int'.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Maps to 'CreatedAt time.Time'. TIMESTAMPTZ (timestamp with time zone) is best for Go's time.Time. Removed UNIQUE.
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Maps to 'UpdatedAt time.Time'. Removed UNIQUE.
    status catalog.status NOT NULL DEFAULT 'draft', -- New field with ENUM type and default
    availability TEXT NOT NULL, -- Maps to 'Availability'. Using snake_case.
    buffer_time INT NOT NULL, -- Maps to 'BufferTime int'. Using snake_case.
    cancellation_hours INT NOT NULL, -- Maps to 'CancellationHours int'. Using snake_case.
    stripe_product_id TEXT NOT NULL, -- Maps to 'StripeProductID string,omitempty'. 'omitempty' means it can be NULL. Removed UNIQUE as descriptions are rarely unique.
    metadata JSONB, -- Maps to 'Metadata map[string]any,omitempty'. JSONB is perfect for storing arbitrary JSON objects and offers efficient querying. 'omitempty' means it can be NULL.

    CONSTRAINT chk_status CHECK (status IN ('published', 'draft', 'archived')),
    CONSTRAINT chk_availability CHECK (availability IN ('online', 'in-person', 'hybrid')),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES catalog.categories(id) ON DELETE RESTRICT
);

-- Now, create the trigger for the catalog.products table
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON catalog.products -- Specify the catalog.products table
FOR EACH ROW EXECUTE FUNCTION update_timestamp(); -- Call the generic update_timestamp function

-- Add comments for better documentation (optional but good practice)
COMMENT ON TABLE catalog.products IS 'Stores information about products available in the system.';
COMMENT ON COLUMN catalog.products.id IS 'Unique identifier for the product.';
COMMENT ON COLUMN catalog.products.name IS 'Name of the product (must be unique).';
COMMENT ON COLUMN catalog.products.description IS 'Detailed description of the product.';
COMMENT ON COLUMN catalog.products.category_id IS 'ID of the category the product belongs to.';
COMMENT ON COLUMN catalog.products.duration IS 'Duration of the product/service in minutes.';
COMMENT ON COLUMN catalog.products.created_at IS 'Timestamp when the product record was created.';
COMMENT ON COLUMN catalog.products.updated_at IS 'Timestamp when the product record was last updated.';
COMMENT ON COLUMN catalog.products.status IS 'Current publication status of the product (e.g., published, draft).';
COMMENT ON COLUMN catalog.products.availability IS 'How the product/service is offered (e.g., online, in-person).';
COMMENT ON COLUMN catalog.products.buffer_time IS 'Buffer time required before/after the product/service, in minutes.';
COMMENT ON COLUMN catalog.products.cancellation_hours IS 'Hours before service start for free cancellation.';
COMMENT ON COLUMN catalog.products.metadata IS 'Flexible JSONB field for additional, product-specific attributes.';

CREATE TABLE IF NOT EXISTS catalog.prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE, -- Foreign key to products table
    stripe_price_id VARCHAR(255) UNIQUE NOT NULL, -- Stores the Stripe Price ID
    amount INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL,
    interval catalog.price_interval NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- To track which prices are currently active locally
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Now, create the trigger for the catalog.products table
CREATE TRIGGER update_prices_updated_at
BEFORE UPDATE ON catalog.prices -- Specify the catalog.products table
FOR EACH ROW EXECUTE FUNCTION update_timestamp(); -- Call the generic update_timestamp function

CREATE TABLE catalog.images (
    id UUID PRIMARY KEY,
    parent_id UUID NOT NULL,
    parent_type VARCHAR(20) NOT NULL CHECK (parent_type IN ('product', 'category')),
    title TEXT,
    s3_key VARCHAR(255) UNIQUE NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add an index for fast lookups by parent
CREATE INDEX parent_lookup_idx ON catalog.images (parent_id, parent_type);
-- Now, create the trigger for the catalog.products table
CREATE TRIGGER update_images_updated_at
BEFORE UPDATE ON catalog.images -- Specify the catalog.products table
FOR EACH ROW EXECUTE FUNCTION update_timestamp(); -- Call the generic update_timestamp function
-- This index guarantees that for any given (parent_id, parent_type) pair, there can be at most one row where the is_active column is true
CREATE UNIQUE INDEX unique_active_image_idx
ON catalog.images (parent_id, parent_type)
WHERE is_active = true;

-- NOTE: could add that constraint at the database level but my domain service already handles it
-- ALTER TABLE catalog.images
-- ADD CONSTRAINT fk_parent_id_exists
-- CHECK (
--     (parent_type = 'product' AND EXISTS (SELECT 1 FROM catalog.products WHERE id = parent_id))
--     OR
--     (parent_type = 'category' AND EXISTS (SELECT 1 FROM catalog.categories WHERE id = parent_id))
-- );

-- Create coupons table
CREATE TABLE IF NOT EXISTS catalog.coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_coupon_id VARCHAR(255) UNIQUE NOT NULL,           -- Stripe Coupon ID
    name TEXT NOT NULL,                                       -- Human-readable name
    percent_off DECIMAL(5,2) CHECK (percent_off > 0 AND percent_off <= 100), -- 25.50 for 25.5%
    amount_off INTEGER CHECK (amount_off > 0),                -- Amount in cents
    currency VARCHAR(3) CHECK (currency ~ '^[A-Z]{3}$'),     -- ISO currency code (3 uppercase letters)
    duration catalog.coupon_duration NOT NULL,              -- How long discount applies
    duration_in_months INTEGER CHECK (duration_in_months > 0), -- Required if duration = 'repeating'
    max_redemptions INTEGER CHECK (max_redemptions > 0),     -- Max usage limit
    times_redeemed INTEGER NOT NULL DEFAULT 0 CHECK (times_redeemed >= 0), -- Current usage count
    valid BOOLEAN NOT NULL DEFAULT TRUE,                     -- Active status
    redeem_by TIMESTAMPTZ,                                   -- Expiry date
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,                                          -- Additional attributes
    
    -- Business logic constraints
    CONSTRAINT chk_discount_type CHECK (
        (percent_off IS NOT NULL AND amount_off IS NULL AND currency IS NULL) OR
        (percent_off IS NULL AND amount_off IS NOT NULL AND currency IS NOT NULL)
    ),
    CONSTRAINT chk_duration_months CHECK (
        (duration = 'repeating' AND duration_in_months IS NOT NULL) OR
        (duration != 'repeating' AND duration_in_months IS NULL)
    ),
    CONSTRAINT chk_times_redeemed_limit CHECK (
        max_redemptions IS NULL OR times_redeemed <= max_redemptions
    )
);

-- Create promotion_codes table
CREATE TABLE IF NOT EXISTS catalog.promotion_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_promotion_id VARCHAR(255) UNIQUE NOT NULL,        -- Stripe Promotion Code ID
    coupon_id UUID NOT NULL,                                 -- FK to coupons
    code VARCHAR(50) NOT NULL UNIQUE,                        -- Customer-facing code
    active BOOLEAN NOT NULL DEFAULT TRUE,                    -- Active status
    max_redemptions INTEGER CHECK (max_redemptions > 0),     -- Max usage for this code
    times_redeemed INTEGER NOT NULL DEFAULT 0 CHECK (times_redeemed >= 0), -- Current usage count
    expires_at TIMESTAMPTZ,                                  -- Code expiry
    first_time_transaction BOOLEAN NOT NULL DEFAULT FALSE,   -- New customer only
    minimum_amount INTEGER CHECK (minimum_amount > 0),       -- Min order amount (cents)
    minimum_amount_currency VARCHAR(3) CHECK (minimum_amount_currency ~ '^[A-Z]{3}$'), -- Currency for min amount
    restrictions JSONB,                                      -- Additional restrictions (currency_options, etc.)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,                                          -- Additional attributes
    
    -- Foreign Key
    CONSTRAINT fk_coupon FOREIGN KEY (coupon_id) REFERENCES catalog.coupons(id) ON DELETE RESTRICT,
    
    -- Business logic constraints
    CONSTRAINT chk_code_format CHECK (code ~ '^[A-Z0-9_-]{3,50}$'), -- Alphanumeric, underscore, hyphen, 3-50 chars
    CONSTRAINT chk_times_redeemed_promo_limit CHECK (
        max_redemptions IS NULL OR times_redeemed <= max_redemptions
    ),
    CONSTRAINT chk_minimum_currency CHECK (
        (minimum_amount IS NOT NULL AND minimum_amount_currency IS NOT NULL) OR
        (minimum_amount IS NULL AND minimum_amount_currency IS NULL)
    )
);

-- Add triggers for updated_at timestamps using existing function
CREATE TRIGGER update_coupons_updated_at
BEFORE UPDATE ON catalog.coupons
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_promotion_codes_updated_at
BEFORE UPDATE ON catalog.promotion_codes
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Add performance indexes
CREATE INDEX idx_coupons_stripe_id ON catalog.coupons(stripe_coupon_id);
CREATE INDEX idx_coupons_valid ON catalog.coupons(valid) WHERE valid = true;
CREATE INDEX idx_coupons_duration ON catalog.coupons(duration);

CREATE INDEX idx_promotion_codes_coupon_id ON catalog.promotion_codes(coupon_id);
CREATE INDEX idx_promotion_codes_code ON catalog.promotion_codes(code);
CREATE INDEX idx_promotion_codes_stripe_id ON catalog.promotion_codes(stripe_promotion_id);
CREATE INDEX idx_promotion_codes_active ON catalog.promotion_codes(active) WHERE active = true;

-- Add table and column comments for coupons
COMMENT ON TABLE catalog.coupons IS 'Stripe coupons defining discount rules, limits, and validity periods';
COMMENT ON COLUMN catalog.coupons.stripe_coupon_id IS 'Unique identifier from Stripe for this coupon';
COMMENT ON COLUMN catalog.coupons.name IS 'Human-readable name for the coupon';
COMMENT ON COLUMN catalog.coupons.percent_off IS 'Percentage discount (mutually exclusive with amount_off)';
COMMENT ON COLUMN catalog.coupons.amount_off IS 'Fixed amount discount in cents (mutually exclusive with percent_off)';
COMMENT ON COLUMN catalog.coupons.currency IS 'ISO currency code, required when amount_off is used';
COMMENT ON COLUMN catalog.coupons.duration IS 'How long the discount applies: once, repeating, or forever';
COMMENT ON COLUMN catalog.coupons.duration_in_months IS 'Number of months discount applies (required when duration = repeating)';
COMMENT ON COLUMN catalog.coupons.max_redemptions IS 'Maximum number of times this coupon can be redeemed globally';
COMMENT ON COLUMN catalog.coupons.times_redeemed IS 'Current number of times this coupon has been redeemed';
COMMENT ON COLUMN catalog.coupons.valid IS 'Whether this coupon is currently valid and can be used';
COMMENT ON COLUMN catalog.coupons.redeem_by IS 'Last date/time this coupon can be redeemed';
COMMENT ON COLUMN catalog.coupons.metadata IS 'Additional metadata for the coupon';

-- Add table and column comments for promotion codes
COMMENT ON TABLE catalog.promotion_codes IS 'Customer-facing promotion codes that reference coupons with additional restrictions';
COMMENT ON COLUMN catalog.promotion_codes.stripe_promotion_id IS 'Unique identifier from Stripe for this promotion code';
COMMENT ON COLUMN catalog.promotion_codes.coupon_id IS 'Foreign key reference to the associated coupon';
COMMENT ON COLUMN catalog.promotion_codes.code IS 'The actual code customers enter (unique, uppercase alphanumeric)';
COMMENT ON COLUMN catalog.promotion_codes.active IS 'Whether this promotion code is currently active';
COMMENT ON COLUMN catalog.promotion_codes.max_redemptions IS 'Maximum number of times this specific code can be used';
COMMENT ON COLUMN catalog.promotion_codes.times_redeemed IS 'Number of times this promotion code has been used';
COMMENT ON COLUMN catalog.promotion_codes.expires_at IS 'Date/time when this promotion code expires';
COMMENT ON COLUMN catalog.promotion_codes.first_time_transaction IS 'Whether code can only be used by new customers';
COMMENT ON COLUMN catalog.promotion_codes.minimum_amount IS 'Minimum order amount in cents required to use this code';
COMMENT ON COLUMN catalog.promotion_codes.minimum_amount_currency IS 'Currency for the minimum amount requirement';
COMMENT ON COLUMN catalog.promotion_codes.restrictions IS 'Additional restrictions (e.g., currency_options) in JSONB format';
COMMENT ON COLUMN catalog.promotion_codes.metadata IS 'Additional metadata for the promotion code';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS catalog.promotion_codes;
DROP TABLE IF EXISTS catalog.coupons;
DROP TABLE IF EXISTS catalog.images;
DROP TABLE IF EXISTS catalog.prices;
DROP TABLE IF EXISTS catalog.products;
DROP TABLE IF EXISTS catalog.categories;
DROP TYPE IF EXISTS catalog.coupon_duration;
DROP TYPE IF EXISTS catalog.price_interval;
DROP TYPE IF EXISTS catalog.status;
DROP SCHEMA IF EXISTS catalog;
-- +goose StatementEnd

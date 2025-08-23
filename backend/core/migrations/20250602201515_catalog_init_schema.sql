-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS products;

CREATE TYPE products.status AS ENUM ('draft', 'published', 'archived');
CREATE TYPE products.price_interval AS ENUM ('month', 'year', 'one_time');

CREATE TABLE IF NOT EXISTS products.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- <-- Change from TEXT to UUID
    name TEXT NOT NULL UNIQUE, -- Maps to 'Name string'. Name should be unique.
    description TEXT, -- Maps to 'Description string,omitempty'. 'omitempty' implies NULLable. Removed UNIQUE as descriptions are not typically unique.
    status products.status NOT NULL DEFAULT 'draft', -- New field with ENUM type and default
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
BEFORE UPDATE ON products.categories
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Add comments for better documentation (optional but good practice)
COMMENT ON TABLE products.categories IS 'Defines categories for products, allowing for flexible grouping and custom attributes.';
COMMENT ON COLUMN products.categories.id IS 'Unique identifier for the category (e.g., UUID).';
COMMENT ON COLUMN products.categories.name IS 'Name of the category (must be unique).';
COMMENT ON COLUMN products.categories.description IS 'Description of the category.';
COMMENT ON COLUMN products.categories.metadata IS 'Flexible JSONB field for additional, category-specific attributes or configurations.';

CREATE TABLE IF NOT EXISTS products.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE, -- Maps to 'Name string'. Keeping UNIQUE as product names are often unique.
    description TEXT NOT NULL, -- Maps to 'Description string,omitempty'. 'omitempty' means it can be NULL. Removed UNIQUE as descriptions are rarely unique.
    category_id UUID NOT NULL, -- Maps to 'CategoryID string'. Using snake_case. Removed UNIQUE, as many products can share a category.
    duration INT NOT NULL, -- Maps to 'Duration int'.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Maps to 'CreatedAt time.Time'. TIMESTAMPTZ (timestamp with time zone) is best for Go's time.Time. Removed UNIQUE.
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Maps to 'UpdatedAt time.Time'. Removed UNIQUE.
    status products.status NOT NULL DEFAULT 'draft', -- New field with ENUM type and default
    availability TEXT NOT NULL, -- Maps to 'Availability'. Using snake_case.
    buffer_time INT NOT NULL, -- Maps to 'BufferTime int'. Using snake_case.
    cancellation_hours INT NOT NULL, -- Maps to 'CancellationHours int'. Using snake_case.
    stripe_product_id TEXT NOT NULL, -- Maps to 'StripeProductID string,omitempty'. 'omitempty' means it can be NULL. Removed UNIQUE as descriptions are rarely unique.
    metadata JSONB, -- Maps to 'Metadata map[string]any,omitempty'. JSONB is perfect for storing arbitrary JSON objects and offers efficient querying. 'omitempty' means it can be NULL.

    CONSTRAINT chk_status CHECK (status IN ('published', 'draft', 'archived')),
    CONSTRAINT chk_availability CHECK (availability IN ('online', 'in-person', 'hybrid')),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES products.categories(id) ON DELETE RESTRICT
);

-- Now, create the trigger for the products.products table
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products.products -- Specify the products.products table
FOR EACH ROW EXECUTE FUNCTION update_timestamp(); -- Call the generic update_timestamp function

-- Add comments for better documentation (optional but good practice)
COMMENT ON TABLE products.products IS 'Stores information about products available in the system.';
COMMENT ON COLUMN products.products.id IS 'Unique identifier for the product.';
COMMENT ON COLUMN products.products.name IS 'Name of the product (must be unique).';
COMMENT ON COLUMN products.products.description IS 'Detailed description of the product.';
COMMENT ON COLUMN products.products.category_id IS 'ID of the category the product belongs to.';
COMMENT ON COLUMN products.products.duration IS 'Duration of the product/service in minutes.';
COMMENT ON COLUMN products.products.created_at IS 'Timestamp when the product record was created.';
COMMENT ON COLUMN products.products.updated_at IS 'Timestamp when the product record was last updated.';
COMMENT ON COLUMN products.products.status IS 'Current publication status of the product (e.g., published, draft).';
COMMENT ON COLUMN products.products.availability IS 'How the product/service is offered (e.g., online, in-person).';
COMMENT ON COLUMN products.products.buffer_time IS 'Buffer time required before/after the product/service, in minutes.';
COMMENT ON COLUMN products.products.cancellation_hours IS 'Hours before service start for free cancellation.';
COMMENT ON COLUMN products.products.metadata IS 'Flexible JSONB field for additional, product-specific attributes.';

CREATE TABLE IF NOT EXISTS products.prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products.products(id) ON DELETE CASCADE, -- Foreign key to products table
    stripe_price_id VARCHAR(255) UNIQUE NOT NULL, -- Stores the Stripe Price ID
    amount INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL,
    interval products.price_interval NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- To track which prices are currently active locally
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Now, create the trigger for the products.products table
CREATE TRIGGER update_prices_updated_at
BEFORE UPDATE ON products.prices -- Specify the products.products table
FOR EACH ROW EXECUTE FUNCTION update_timestamp(); -- Call the generic update_timestamp function

CREATE TABLE products.images (
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
CREATE INDEX parent_lookup_idx ON products.images (parent_id, parent_type);
-- Now, create the trigger for the products.products table
CREATE TRIGGER update_images_updated_at
BEFORE UPDATE ON products.images -- Specify the products.products table
FOR EACH ROW EXECUTE FUNCTION update_timestamp(); -- Call the generic update_timestamp function
-- This index guarantees that for any given (parent_id, parent_type) pair, there can be at most one row where the is_active column is true
CREATE UNIQUE INDEX unique_active_image_idx
ON products.images (parent_id, parent_type)
WHERE is_active = true;

-- NOTE: could add that constraint at the database level but my domain service already handles it
-- ALTER TABLE products.images
-- ADD CONSTRAINT fk_parent_id_exists
-- CHECK (
--     (parent_type = 'product' AND EXISTS (SELECT 1 FROM products.products WHERE id = parent_id))
--     OR
--     (parent_type = 'category' AND EXISTS (SELECT 1 FROM products.categories WHERE id = parent_id))
-- );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products.images;
DROP TABLE products.prices;
DROP TABLE products.products;
DROP TABLE products.categories;
DROP SCHEMA IF EXISTS products;
-- +goose StatementEnd

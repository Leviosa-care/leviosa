-- +goose Up
-- +goose StatementBegin

-- Create booking schema
CREATE SCHEMA IF NOT EXISTS booking;

-- Buildings: Physical locations containing treatment rooms
CREATE TABLE booking.buildings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Name and address (encrypted for GDPR compliance)
    name_encrypted BYTEA NOT NULL,
    name_hash TEXT NOT NULL,
    address_encrypted BYTEA NOT NULL,
    address_hash TEXT NOT NULL,
    city_encrypted BYTEA NOT NULL,
    city_hash TEXT NOT NULL,
    postal_code_encrypted BYTEA NOT NULL,
    country_encrypted BYTEA NOT NULL,
    country_hash TEXT NOT NULL,

    -- Business information
    description_encrypted BYTEA,
    phone_encrypted BYTEA,
    email_encrypted BYTEA,

    -- Encryption metadata (required by encx)
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Administrative fields
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    metadata JSONB -- Maps to 'Metadata map[string]any,omitempty'. JSONB for flexible key-value pairs. 'omitempty' implies NULLable.
);

-- Indexes for searchable fields (hash-based filtering)
CREATE INDEX idx_buildings_name_hash ON booking.buildings(name_hash);
CREATE INDEX idx_buildings_address_hash ON booking.buildings(address_hash);
CREATE INDEX idx_buildings_city_hash ON booking.buildings(city_hash);
CREATE INDEX idx_buildings_country_hash ON booking.buildings(country_hash);
CREATE INDEX idx_buildings_is_active ON booking.buildings(is_active);

-- Rooms: Individual treatment spaces within buildings
CREATE TABLE booking.rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    building_id UUID NOT NULL REFERENCES booking.buildings(id) ON DELETE CASCADE,

    -- Room identification (encrypted)
    name_encrypted BYTEA NOT NULL,
    name_hash TEXT NOT NULL,
    description_encrypted BYTEA,
    room_number_encrypted BYTEA,
    room_number_hash TEXT,

    -- Room specifications
    capacity INTEGER NOT NULL DEFAULT 1 CHECK (capacity > 0),
    equipment_encrypted BYTEA, -- JSON array of equipment/amenities

    -- Pricing and availability
    hourly_rate_cents INTEGER, -- Optional base rate in cents

    -- Encryption metadata (required by encx)
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Administrative fields
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    metadata JSONB -- Maps to 'Metadata map[string]any,omitempty'. JSONB for flexible key-value pairs. 'omitempty' implies NULLable.
);

-- Room allocations: Partner assignments to rooms (dedicated vs shared access)
CREATE TABLE booking.room_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES booking.rooms(id) ON DELETE CASCADE,
    partner_id UUID NOT NULL, -- References auth.partners.id (from authuser service)

    -- Allocation type
    allocation_type VARCHAR(20) NOT NULL CHECK (allocation_type IN ('dedicated', 'shared')),

    -- Time-based allocation (optional - for dedicated allocations with time limits)
    start_date DATE,
    end_date DATE,

    -- Administrative fields
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    metadata JSONB, -- Maps to 'Metadata map[string]any,omitempty'. JSONB for flexible key-value pairs. 'omitempty' implies NULLable.

    -- Constraints
    UNIQUE(room_id, partner_id, allocation_type),
    CHECK ((allocation_type = 'dedicated' AND start_date IS NOT NULL) OR allocation_type = 'shared')
);

-- Availability: Time slots partners offer for services
CREATE TABLE booking.availabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL, -- References auth.partners.id (from authuser service)
    room_id UUID NOT NULL REFERENCES booking.rooms(id) ON DELETE CASCADE,

    -- Time slot definition
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Service offering details
    service_type_encrypted BYTEA, -- Type of service offered
    price_cents INTEGER, -- Price override for this specific availability
    max_capacity INTEGER NOT NULL DEFAULT 1 CHECK (max_capacity > 0),

    -- Availability metadata
    notes_encrypted BYTEA, -- Internal notes for the partner
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    recurrence_pattern_encrypted BYTEA, -- JSON with recurrence rules if applicable

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'booked', 'cancelled', 'blocked')),

    -- Encryption metadata (required by encx)
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Administrative fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    metadata JSONB, -- Maps to 'Metadata map[string]any,omitempty'. JSONB for flexible key-value pairs. 'omitempty' implies NULLable.

    -- Constraints
    CHECK (end_time > start_time),
    CHECK (start_time >= NOW()) -- Cannot create availability in the past
);

-- Bookings: Client reservations of partner availability
CREATE TABLE booking.bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    availability_id UUID NOT NULL REFERENCES booking.availabilities(id) ON DELETE RESTRICT,
    client_id UUID NOT NULL, -- References auth.users.id (from authuser service)
    partner_id UUID NOT NULL, -- References auth.partners.id (denormalized for queries)
    room_id UUID NOT NULL, -- References booking.rooms.id (denormalized for queries)

    -- Booking details (encrypted for GDPR compliance)
    client_notes_encrypted BYTEA, -- Special requests from client
    partner_notes_encrypted BYTEA, -- Private notes from partner

    -- Pricing information
    total_price_cents INTEGER NOT NULL CHECK (total_price_cents >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',

    -- Payment tracking
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded')),
    payment_intent_id VARCHAR(255), -- Stripe payment intent ID

    -- Booking lifecycle
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed' CHECK (status IN ('confirmed', 'cancelled', 'completed', 'no_show')),
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancellation_reason_encrypted BYTEA,

    -- Encryption metadata (required by encx)
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Administrative fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    metadata JSONB -- Maps to 'Metadata map[string]any,omitempty'. JSONB for flexible key-value pairs. 'omitempty' implies NULLable.
);

-- Create indexes for performance
CREATE INDEX idx_buildings_active ON booking.buildings(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_rooms_building_active ON booking.rooms(building_id, is_active) WHERE is_active = TRUE;

-- Indexes for searchable fields (hash-based filtering)
CREATE INDEX idx_rooms_name_hash ON booking.rooms(name_hash);
CREATE INDEX idx_rooms_room_number_hash ON booking.rooms(room_number_hash);

CREATE INDEX idx_room_allocations_partner ON booking.room_allocations(partner_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_room_allocations_room ON booking.room_allocations(room_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_availabilities_partner_time ON booking.availabilities(partner_id, start_time, end_time);
CREATE INDEX idx_availabilities_room_time ON booking.availabilities(room_id, start_time, end_time);
CREATE INDEX idx_availabilities_status_time ON booking.availabilities(status, start_time) WHERE status = 'available';
CREATE INDEX idx_bookings_client ON booking.bookings(client_id, created_at DESC);
CREATE INDEX idx_bookings_partner ON booking.bookings(partner_id, created_at DESC);
CREATE INDEX idx_bookings_availability ON booking.bookings(availability_id);
CREATE INDEX idx_bookings_status ON booking.bookings(status, created_at DESC);
CREATE INDEX idx_bookings_payment_status ON booking.bookings(payment_status) WHERE payment_status IN ('pending', 'failed');

-- Create update timestamp triggers
CREATE OR REPLACE FUNCTION booking.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER buildings_update_timestamp BEFORE UPDATE ON booking.buildings
    FOR EACH ROW EXECUTE FUNCTION booking.update_timestamp();

CREATE TRIGGER rooms_update_timestamp BEFORE UPDATE ON booking.rooms
    FOR EACH ROW EXECUTE FUNCTION booking.update_timestamp();

CREATE TRIGGER room_allocations_update_timestamp BEFORE UPDATE ON booking.room_allocations
    FOR EACH ROW EXECUTE FUNCTION booking.update_timestamp();

CREATE TRIGGER availabilities_update_timestamp BEFORE UPDATE ON booking.availabilities
    FOR EACH ROW EXECUTE FUNCTION booking.update_timestamp();

CREATE TRIGGER bookings_update_timestamp BEFORE UPDATE ON booking.bookings
    FOR EACH ROW EXECUTE FUNCTION booking.update_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop schema and all tables (cascades to indexes and triggers)
DROP SCHEMA IF EXISTS booking CASCADE;

-- +goose StatementEnd

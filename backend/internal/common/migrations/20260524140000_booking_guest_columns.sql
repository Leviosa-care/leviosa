-- +goose Up
-- +goose StatementBegin

-- Make client_id nullable to support guest bookings (ADR 0008)
ALTER TABLE booking.bookings ALTER COLUMN client_id DROP NOT NULL;

-- Add encrypted guest contact fields for bookings without an account
ALTER TABLE booking.bookings ADD COLUMN guest_first_name_encrypted BYTEA;
ALTER TABLE booking.bookings ADD COLUMN guest_last_name_encrypted BYTEA;
ALTER TABLE booking.bookings ADD COLUMN guest_email_encrypted BYTEA;
ALTER TABLE booking.bookings ADD COLUMN guest_phone_encrypted BYTEA;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Abort if guest bookings exist; restoring NOT NULL would fail silently otherwise.
DO $$
DECLARE guest_count BIGINT;
BEGIN
  SELECT COUNT(*) INTO guest_count FROM booking.bookings WHERE client_id IS NULL;
  IF guest_count > 0 THEN
    RAISE EXCEPTION 'Cannot roll back migration: % guest booking(s) have client_id IS NULL. Delete or migrate them before rolling back.', guest_count;
  END IF;
END $$;

-- Remove guest contact fields
ALTER TABLE booking.bookings DROP COLUMN IF EXISTS guest_phone_encrypted;
ALTER TABLE booking.bookings DROP COLUMN IF EXISTS guest_email_encrypted;
ALTER TABLE booking.bookings DROP COLUMN IF EXISTS guest_last_name_encrypted;
ALTER TABLE booking.bookings DROP COLUMN IF EXISTS guest_first_name_encrypted;

-- Restore NOT NULL constraint on client_id
ALTER TABLE booking.bookings ALTER COLUMN client_id SET NOT NULL;

-- +goose StatementEnd

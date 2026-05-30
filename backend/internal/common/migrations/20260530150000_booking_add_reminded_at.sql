-- +goose Up
-- +goose StatementBegin

-- Add nullable reminded_at column for the booking reminder scheduler (issue 006).
-- Tracks when a reminder notification was sent for a booking.
ALTER TABLE booking.bookings ADD COLUMN reminded_at TIMESTAMP WITH TIME ZONE;

-- Partial index to keep scheduler queries fast: only indexes bookings that are
-- still eligible for a reminder (confirmed, not yet reminded).
CREATE INDEX idx_bookings_reminder_due ON booking.bookings(reminded_at)
    WHERE reminded_at IS NULL AND status = 'confirmed';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS booking.idx_bookings_reminder_due;
ALTER TABLE booking.bookings DROP COLUMN IF EXISTS reminded_at;

-- +goose StatementEnd

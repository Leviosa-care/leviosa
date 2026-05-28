-- +goose Up
-- +goose StatementBegin

-- Add nullable token column for signed booking tokens (issue 002).
-- Existing bookings created before this feature will have NULL tokens.
ALTER TABLE booking.bookings ADD COLUMN token TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE booking.bookings DROP COLUMN IF EXISTS token;

-- +goose StatementEnd

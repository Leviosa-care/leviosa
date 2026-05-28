package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetGuestBookings retrieves all bookings with client_id IS NULL.
func (r *Repository) GetGuestBookings(ctx context.Context) ([]*domain.BookingEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			b.id, b.availability_id, b.client_id, b.user_id, b.room_id,
			b.product_id_encrypted, b.slot_start_time_encrypted, b.slot_end_time_encrypted,
			b.client_notes_encrypted, b.partner_notes_encrypted,
			b.total_price_cents, b.currency, b.payment_status, b.payment_intent_id,
			b.status, b.cancelled_at, b.cancellation_reason_encrypted, b.completed_at,
			b.guest_first_name_encrypted, b.guest_last_name_encrypted,
			b.guest_email_encrypted, b.guest_phone_encrypted,
			b.token,
			b.created_at, b.updated_at,
			b.dek_encrypted, b.key_version, b.metadata
		FROM %s.bookings b
		WHERE b.client_id IS NULL
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get guest bookings", err)
	}
	defer rows.Close()

	var bookings []*domain.BookingEncx
	for rows.Next() {
		b := &domain.BookingEncx{}
		if err := rows.Scan(
			&b.ID, &b.AvailabilityID, &b.ClientID, &b.PartnerID, &b.RoomID,
			&b.ProductIDEncrypted, &b.SlotStartTimeEncrypted, &b.SlotEndTimeEncrypted,
			&b.ClientNotesEncrypted, &b.PartnerNotesEncrypted,
			&b.TotalPriceCents, &b.Currency, &b.PaymentStatus, &b.PaymentIntentID,
			&b.Status, &b.CancelledAt, &b.CancellationReasonEncrypted, &b.CompletedAt,
			&b.GuestFirstNameEncrypted, &b.GuestLastNameEncrypted,
			&b.GuestEmailEncrypted, &b.GuestPhoneEncrypted,
			&b.Token,
			&b.CreatedAt, &b.UpdatedAt,
			&b.DEKEncrypted, &b.KeyVersion, &b.Metadata,
		); err != nil {
			return nil, errs.ClassifyPgError("scan guest booking", err)
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate guest bookings", err)
	}

	return bookings, nil
}

// SetBookingClientID atomically sets client_id on a single booking, only if
// the booking's client_id is currently NULL. Returns true if the row was updated.
func (r *Repository) SetBookingClientID(ctx context.Context, bookingID, clientID uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		UPDATE %s.bookings SET client_id = $1, updated_at = now()
		WHERE id = $2 AND client_id IS NULL
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, clientID, bookingID)
	if err != nil {
		return false, errs.ClassifyPgError("set booking client_id", err)
	}

	return result.RowsAffected() > 0, nil
}

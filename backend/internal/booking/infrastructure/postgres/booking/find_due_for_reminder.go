package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// FindBookingsDueForReminder returns all confirmed bookings where reminded_at IS NULL.
// Time filtering on slot_start_time is performed at the service layer after
// decryption because slot_start_time is stored encrypted for GDPR compliance.
func (r *Repository) FindBookingsDueForReminder(ctx context.Context) ([]*domain.BookingEncx, error) {
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
			b.reminded_at,
			b.created_at, b.updated_at,
			b.dek_encrypted, b.key_version, b.metadata
		FROM %s.bookings b
		WHERE b.status = 'confirmed' AND b.reminded_at IS NULL
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("find bookings due for reminder", err)
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
			&b.RemindedAt,
			&b.CreatedAt, &b.UpdatedAt,
			&b.DEKEncrypted, &b.KeyVersion, &b.Metadata,
		); err != nil {
			return nil, errs.ClassifyPgError("scan booking due for reminder", err)
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate bookings due for reminder", err)
	}

	return bookings, nil
}

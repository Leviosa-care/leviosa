package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByAvailabilityID(ctx context.Context, availabilityID uuid.UUID) (*domain.BookingEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, availability_id, client_id, user_id, room_id,
			product_id_encrypted, slot_start_time_encrypted, slot_end_time_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			guest_first_name_encrypted, guest_last_name_encrypted,
			guest_email_encrypted, guest_phone_encrypted,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.bookings
		WHERE availability_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, r.schema)

	bookingEncx := &domain.BookingEncx{}
	err := r.pool.QueryRow(ctx, query, availabilityID).Scan(
		&bookingEncx.ID,
		&bookingEncx.AvailabilityID,
		&bookingEncx.ClientID,
		&bookingEncx.PartnerID,
		&bookingEncx.RoomID,
		&bookingEncx.ProductIDEncrypted,
		&bookingEncx.SlotStartTimeEncrypted,
		&bookingEncx.SlotEndTimeEncrypted,
		&bookingEncx.ClientNotesEncrypted,
		&bookingEncx.PartnerNotesEncrypted,
		&bookingEncx.TotalPriceCents,
		&bookingEncx.Currency,
		&bookingEncx.PaymentStatus,
		&bookingEncx.PaymentIntentID,
		&bookingEncx.Status,
		&bookingEncx.CancelledAt,
		&bookingEncx.CancellationReasonEncrypted,
		&bookingEncx.CompletedAt,
		&bookingEncx.GuestFirstNameEncrypted,
		&bookingEncx.GuestLastNameEncrypted,
		&bookingEncx.GuestEmailEncrypted,
		&bookingEncx.GuestPhoneEncrypted,
		&bookingEncx.CreatedAt,
		&bookingEncx.UpdatedAt,
		&bookingEncx.DEKEncrypted,
		&bookingEncx.KeyVersion,
		&bookingEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get booking by availability id", err)
	}

	return bookingEncx, nil
}
package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, booking *domain.BookingEncx) error {
	query := fmt.Sprintf(`
	UPDATE %s.bookings SET
			availability_id = $2,
			client_id = $3,
			user_id = $4,
			room_id = $5,
			product_id_encrypted = $6,
			slot_start_time_encrypted = $7,
			slot_end_time_encrypted = $8,
			client_notes_encrypted = $9,
			partner_notes_encrypted = $10,
			total_price_cents = $11,
			currency = $12,
			payment_status = $13,
			payment_intent_id = $14,
			status = $15,
			cancelled_at = $16,
			cancellation_reason_encrypted = $17,
			completed_at = $18,
			guest_first_name_encrypted = $19,
			guest_last_name_encrypted = $20,
			guest_email_encrypted = $21,
			guest_phone_encrypted = $22,
			token = $23,
			reminded_at = $24,
			updated_at = $25,
			dek_encrypted = $26,
			key_version = $27,
			metadata = $28
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		booking.ID,
		booking.AvailabilityID,
		booking.ClientID,
		booking.PartnerID,
		booking.RoomID,
		booking.ProductIDEncrypted,
		booking.SlotStartTimeEncrypted,
		booking.SlotEndTimeEncrypted,
		booking.ClientNotesEncrypted,
		booking.PartnerNotesEncrypted,
		booking.TotalPriceCents,
		booking.Currency,
		booking.PaymentStatus,
		booking.PaymentIntentID,
		booking.Status,
		booking.CancelledAt,
		booking.CancellationReasonEncrypted,
		booking.CompletedAt,
		booking.GuestFirstNameEncrypted,
		booking.GuestLastNameEncrypted,
		booking.GuestEmailEncrypted,
		booking.GuestPhoneEncrypted,
		booking.Token,
		booking.RemindedAt,
		booking.UpdatedAt,
		booking.DEKEncrypted,
		booking.KeyVersion,
		booking.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("update booking", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}

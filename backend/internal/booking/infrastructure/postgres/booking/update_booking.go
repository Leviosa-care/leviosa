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
			partner_id = $4,
			room_id = $5,
			productid_encrypted = $6,
			slotstarttime_encrypted = $7,
			slotendtime_encrypted = $8,
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
			updated_at = $19,
			dek_encrypted = $20,
			key_version = $21,
			metadata = $22
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


package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, booking *domain.BookingEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.bookings (
			id, availability_id, client_id, partner_id, room_id,
			productid_encrypted, slotstarttime_encrypted, slotendtime_encrypted,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted, completed_at,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
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
		booking.CreatedAt,
		booking.UpdatedAt,
		booking.DEKEncrypted,
		booking.KeyVersion,
		booking.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create booking", err)
	}

	return nil
}

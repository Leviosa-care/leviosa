package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, booking *domain.Booking) error {
	// Encrypt sensitive fields
	if err := r.crypto.EncryptStruct(ctx, booking); err != nil {
		return fmt.Errorf("encrypt booking data: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.bookings (
			id, availability_id, client_id, partner_id, room_id,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		booking.ID,
		booking.AvailabilityID,
		booking.ClientID,
		booking.PartnerID,
		booking.RoomID,
		booking.ClientNotesEncrypted,
		booking.PartnerNotesEncrypted,
		booking.TotalPriceCents,
		booking.Currency,
		booking.PaymentStatus,
		booking.PaymentIntentID,
		booking.Status,
		booking.CancelledAt,
		booking.CancellationReasonEncrypted,
		booking.CreatedAt,
		booking.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create booking", err)
	}

	return nil
}
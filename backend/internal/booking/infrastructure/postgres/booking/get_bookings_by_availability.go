package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetBookingsByAvailability retrieves all bookings for a specific availability.
// Used for slot overlap detection in the slot-based booking system.
// Returns bookings sorted by slot_start_time ASC.
func (r *Repository) GetBookingsByAvailability(ctx context.Context, availabilityID uuid.UUID) ([]*domain.Booking, error) {
	query := `
		SELECT
			id,
			availabilityid,
			clientid,
			partnerid,
			roomid,
			totalpricecents,
			currency,
			paymentstatus,
			paymentintentid,
			status,
			cancelledat,
			completedat,
			createdat,
			updatedat,
			-- Encrypted fields
			productid_encrypted,
			slotstarttime_encrypted,
			slotendtime_encrypted,
			clientnotes_encrypted,
			partnernotes_encrypted,
			cancellationreason_encrypted,
			-- Encryption metadata
			dek_encrypted,
			key_version,
			metadata
		FROM booking.bookings
		WHERE availabilityid = $1
		ORDER BY slotstarttime_encrypted ASC
	`

	rows, err := r.pool.Query(ctx, query, availabilityID)
	if err != nil {
		return nil, errs.ClassifyPgError("get bookings by availability", err)
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		bookingEncx := &domain.BookingEncx{}

		err := rows.Scan(
			&bookingEncx.ID,
			&bookingEncx.AvailabilityID,
			&bookingEncx.ClientID,
			&bookingEncx.PartnerID,
			&bookingEncx.RoomID,
			&bookingEncx.TotalPriceCents,
			&bookingEncx.Currency,
			&bookingEncx.PaymentStatus,
			&bookingEncx.PaymentIntentID,
			&bookingEncx.Status,
			&bookingEncx.CancelledAt,
			&bookingEncx.CompletedAt,
			&bookingEncx.CreatedAt,
			&bookingEncx.UpdatedAt,
			// Encrypted fields
			&bookingEncx.ProductIDEncrypted,
			&bookingEncx.SlotStartTimeEncrypted,
			&bookingEncx.SlotEndTimeEncrypted,
			&bookingEncx.ClientNotesEncrypted,
			&bookingEncx.PartnerNotesEncrypted,
			&bookingEncx.CancellationReasonEncrypted,
			// Encryption metadata
			&bookingEncx.DEKEncrypted,
			&bookingEncx.KeyVersion,
			&bookingEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan booking row", err)
		}

		// Decrypt booking
		booking, err := domain.DecryptBookingEncx(ctx, r.crypto, bookingEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt booking %s: %w", bookingEncx.ID, err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate booking rows", err)
	}

	// Return empty slice instead of nil if no bookings found
	if bookings == nil {
		return []*domain.Booking{}, nil
	}

	return bookings, nil
}

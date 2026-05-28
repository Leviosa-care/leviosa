package booking

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
)

// ClaimBookings links all guest bookings whose decrypted GuestEmail matches the
// provided email to the given client. Bookings that already have an owner are
// skipped. Returns the number of bookings claimed.
func (s *BookingService) ClaimBookings(ctx context.Context, clientID uuid.UUID, email string) (int, error) {
	if clientID == uuid.Nil {
		return 0, fmt.Errorf("client ID is required")
	}
	if strings.TrimSpace(email) == "" {
		return 0, fmt.Errorf("email is required")
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(email))

	// Fetch all unowned bookings from the repository
	guestBookings, err := s.bookingRepo.GetGuestBookings(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch guest bookings: %w", err)
	}

	var claimed int
	for _, encx := range guestBookings {
		// Decrypt to access GuestEmail
		decrypted, decErr := domain.DecryptBookingEncx(ctx, s.crypto, encx)
		if decErr != nil {
			slog.WarnContext(ctx, "claim: skipping booking, decryption failed",
				"booking_id", encx.ID, "err", decErr)
			continue
		}

		guestEmail := strings.TrimSpace(strings.ToLower(decrypted.GuestEmail))
		if guestEmail != normalizedEmail {
			continue
		}

		// Atomically set client_id (only if still NULL)
		updated, setErr := s.bookingRepo.SetBookingClientID(ctx, encx.ID, clientID)
		if setErr != nil {
			slog.WarnContext(ctx, "claim: failed to set client_id",
				"booking_id", encx.ID, "client_id", clientID, "err", setErr)
			continue
		}
		if updated {
			claimed++
		}
	}

	if claimed > 0 {
		slog.InfoContext(ctx, "claim: linked guest bookings to client",
			"client_id", clientID, "count", claimed)
	}

	return claimed, nil
}

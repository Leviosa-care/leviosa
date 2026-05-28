package booking

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// LookupBookingByToken verifies a booking token and returns the booking
// converted to a public response DTO. Returns an error if the token is
// invalid, expired, or the booking does not exist.
func (s *BookingService) LookupBookingByToken(ctx context.Context, token string) (*domain.PublicBookingLookupResponse, error) {
	if len(s.tokenSecret) == 0 {
		return nil, fmt.Errorf("booking token feature is not configured")
	}

	bookingID, err := domain.VerifyBookingToken(token, s.tokenSecret)
	if err != nil {
		if domain.IsBookingTokenError(err) {
			return nil, errs.NewUnauthorizedErr(err.Error())
		}
		return nil, fmt.Errorf("verify booking token: %w", err)
	}

	booking, err := s.GetBooking(ctx, bookingID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewUnauthorizedErr("booking not found")
		}
		return nil, fmt.Errorf("get booking for lookup: %w", err)
	}

	return s.toPublicResponse(ctx, booking), nil
}

// LookupBookingByRefAndContact looks up a guest booking by reference ID
// combined with either an email or phone number. The backend decrypts
// guest contact fields server-side for comparison — plaintext guest data
// is never returned in the response.
func (s *BookingService) LookupBookingByRefAndContact(ctx context.Context, ref string, email, phone string) (*domain.PublicBookingLookupResponse, error) {
	bookingID, err := uuid.Parse(ref)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid booking reference format")
	}

	booking, err := s.GetBooking(ctx, bookingID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewUnauthorizedErr("booking not found")
		}
		return nil, fmt.Errorf("get booking for lookup: %w", err)
	}

	// Only guest bookings can be looked up via ref+contact
	if !booking.IsGuestBooking() {
		return nil, errs.NewUnauthorizedErr("this booking is linked to a registered account")
	}

	// Compare contact fields (server-side decryption already done by GetBooking)
	email = strings.TrimSpace(strings.ToLower(email))
	phone = strings.TrimSpace(phone)

	guestEmail := strings.TrimSpace(strings.ToLower(booking.GuestEmail))
	guestPhone := strings.TrimSpace(booking.GuestPhone)

	matched := false
	if email != "" && guestEmail != "" && email == guestEmail {
		matched = true
	}
	if phone != "" && guestPhone != "" && phone == guestPhone {
		matched = true
	}

	if !matched {
		return nil, errs.NewUnauthorizedErr("contact information does not match")
	}

	return s.toPublicResponse(ctx, booking), nil
}

// toPublicResponse converts a decrypted booking to a public lookup response.
// No guest contact fields or internal IDs are exposed.
func (s *BookingService) toPublicResponse(ctx context.Context, booking *domain.Booking) *domain.PublicBookingLookupResponse {
	resp := &domain.PublicBookingLookupResponse{
		ID:              booking.ID,
		SlotStartTime:   booking.SlotStartTime,
		SlotEndTime:     booking.SlotEndTime,
		Status:          booking.Status,
		TotalPriceCents: booking.TotalPriceCents,
		Currency:        booking.Currency,
		PaymentStatus:   booking.PaymentStatus,
		ProductName:     s.resolveProductName(ctx, booking.ProductID),
		PartnerName:     s.resolveUserName(ctx, booking.PartnerID, "Praticien inconnu"),
	}

	return resp
}

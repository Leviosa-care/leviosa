package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *BookingService) CreateBooking(ctx context.Context, availabilityID, clientID uuid.UUID, clientNotes string) (*domain.Booking, error) {
	// Get availability and verify it's bookable
	availabilityEncx, err := s.availabilityRepo.GetByID(ctx, availabilityID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("availability not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("get availability: %w", err)
	}

	availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("availability", err)
	}

	// Check if availability is still available for booking
	if !availability.IsAvailable() {
		return nil, fmt.Errorf("availability is no longer available for booking")
	}

	// Check if there's already a booking for this availability
	existingBooking, err := s.bookingRepo.GetByAvailabilityID(ctx, availabilityID)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("check existing booking: %w", err)
	}

	if existingBooking != nil {
		return nil, fmt.Errorf("availability is already booked")
	}

	// Calculate price (use availability price or default to 0 for now)
	totalPriceCents := 0
	if availability.PriceCents != nil {
		totalPriceCents = *availability.PriceCents
	}

	// Create booking entity
	booking, err := domain.NewBooking(
		availabilityID,
		clientID,
		availability.UserID,
		availability.RoomID,
		totalPriceCents,
		"EUR", // Default currency
	)
	if err != nil {
		return nil, fmt.Errorf("create booking entity: %w", err)
	}

	// Set client notes if provided
	if clientNotes != "" {
		booking.SetClientNotes(clientNotes)
	}

	// Create payment intent if booking has a price
	if totalPriceCents > 0 {
		description := fmt.Sprintf("Booking for availability %s", availabilityID.String())
		metadata := map[string]string{
			"booking_id":      booking.ID.String(),
			"availability_id": availabilityID.String(),
			"client_id":       clientID.String(),
			"partner_id":      availability.UserID.String(),
		}

		paymentIntentID, _, err := s.paymentService.CreatePaymentIntent(
			ctx,
			totalPriceCents,
			booking.Currency,
			description,
			metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("create payment intent: %w", err)
		}

		// Update booking with payment intent ID
		booking.SetPaymentIntentID(paymentIntentID)
	}

	// Mark availability as booked (atomic operation would be ideal here)
	if err := availability.MarkAsBooked(); err != nil {
		return nil, fmt.Errorf("mark availability as booked: %w", err)
	}

	availabilityEncx, err = domain.ProcessAvailabilityEncx(ctx, s.crypto, availability)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("availability", err)
	}

	// Update availability in repository
	if err := s.availabilityRepo.Update(ctx, availabilityEncx); err != nil {
		return nil, fmt.Errorf("update availability status: %w", err)
	}

	// Persist booking to repository
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		// Rollback availability status on failure
		availability.MarkAsAvailable()
		if err := s.availabilityRepo.Update(ctx, availabilityEncx); err != nil {

		}
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return booking, nil
}


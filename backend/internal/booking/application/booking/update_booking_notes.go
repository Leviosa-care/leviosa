package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *BookingService) UpdateBookingNotes(ctx context.Context, id uuid.UUID, clientNotes, partnerNotes string) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for notes update: %w", err)
	}

	// Update notes
	if clientNotes != "" {
		booking.SetClientNotes(clientNotes)
	}
	if partnerNotes != "" {
		booking.SetPartnerNotes(partnerNotes)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update booking notes: %w", err)
	}

	return booking, nil
}
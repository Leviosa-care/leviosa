package availability

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *AvailabilityService) UpdateAvailability(ctx context.Context, id uuid.UUID, startTime, endTime time.Time, serviceType string, priceCents *int, notes string) (*domain.Availability, error) {
	// Get existing availability
	availabilityEncx, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get availability by ID for update: %w", err)
	}

	availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("availability", err)
	}

	// Check if availability can be updated (not booked)
	if availability.Status == domain.AvailabilityStatusBooked {
		return nil, fmt.Errorf("cannot update booked availability")
	}

	// If time is being changed, check for conflicts
	if !startTime.Equal(availability.StartTime) || !endTime.Equal(availability.EndTime) {
		// Check partner still has room access for new time
		hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, availability.UserID, availability.RoomID, startTime)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				return nil, fmt.Errorf("partner does not have allocation for this room at new time")
			}
			return nil, fmt.Errorf("check partner room access: %w", err)
		}

		if !hasAccess.IsActiveAt(startTime) || !hasAccess.IsActiveAt(endTime) {
			return nil, fmt.Errorf("partner does not have access to room during new time period")
		}

		// Check for scheduling conflicts (excluding this availability)
		hasConflict, err := s.availabilityRepo.CheckConflict(ctx, availability.UserID, startTime, endTime, &id)
		if err != nil {
			return nil, fmt.Errorf("check availability conflict: %w", err)
		}

		if hasConflict {
			return nil, fmt.Errorf("updated time conflicts with existing availability")
		}

		// Update time slot
		if err := availability.UpdateTimeSlot(startTime, endTime); err != nil {
			return nil, fmt.Errorf("update time slot: %w", err)
		}
	}

	// Update service details
	availability.SetServiceDetails(serviceType, priceCents, notes)

	updatedAvailabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, s.crypto, availability)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("availability", err)
	}

	// Persist changes
	if err := s.availabilityRepo.Update(ctx, updatedAvailabilityEncx); err != nil {
		return nil, fmt.Errorf("update availability: %w", err)
	}

	return availability, nil
}

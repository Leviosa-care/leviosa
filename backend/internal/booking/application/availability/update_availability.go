package availability

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/hengadev/encx"
)

// func (s *AvailabilityService) UpdateAvailability(ctx context.Context, id uuid.UUID, startTime, endTime time.Time, serviceType string, priceCents *int, notes string) (*domain.Availability, error) {
func (s *AvailabilityService) UpdateAvailability(ctx context.Context, request *domain.UpdateAvailabilityRequest) (*domain.Availability, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing availability
	availabilityEncx, err := s.availabilityRepo.GetByID(ctx, request.ID)
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
	if (request.StartTime != nil && !request.StartTime.Equal(availability.StartTime)) ||
		(request.EndTime != nil && !request.EndTime.Equal(availability.EndTime)) {

		// Use new times if provided, otherwise keep existing
		newStartTime := availability.StartTime
		newEndTime := availability.EndTime
		if request.StartTime != nil {
			newStartTime = *request.StartTime
		}
		if request.EndTime != nil {
			newEndTime = *request.EndTime
		}

		// Check partner still has room access for new time
		userIDBytes, err := encx.SerializeValue(availability.UserID)
		if err != nil {
			return nil, fmt.Errorf("serialize user ID for hashing: %w", err)
		}
		userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

		roomAllocationEncx, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, availability.RoomID, newStartTime)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				return nil, fmt.Errorf("partner does not have allocation for this room at new time")
			}
			return nil, fmt.Errorf("check partner room access: %w", err)
		}

		hasAccess, err := domain.DecryptRoomAllocationEncx(ctx, s.crypto, roomAllocationEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("room allocation", err)
		}

		if !hasAccess.IsActiveAt(newStartTime) || !hasAccess.IsActiveAt(newEndTime) {
			return nil, fmt.Errorf("partner does not have access to room during new time period")
		}

		// Check for scheduling conflicts (excluding this availability)
		hasConflict, err := s.availabilityRepo.CheckConflict(ctx, availability.UserID, newStartTime, newEndTime, &request.ID)
		if err != nil {
			return nil, fmt.Errorf("check availability conflict: %w", err)
		}

		if hasConflict {
			return nil, fmt.Errorf("updated time conflicts with existing availability")
		}

		// Update time slot
		if err := availability.UpdateTimeSlot(newStartTime, newEndTime); err != nil {
			return nil, fmt.Errorf("update time slot: %w", err)
		}
	}

	// Update service details (handle optional pointer fields)
	var serviceType string
	var notes string
	if request.ServiceType != nil {
		serviceType = *request.ServiceType
	}
	if request.Notes != nil {
		notes = *request.Notes
	}
	availability.SetServiceDetails(serviceType, request.PriceCents, notes)

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

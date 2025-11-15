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

func (s *AvailabilityService) CreateAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int) (*domain.Availability, error) {
	// Validate partner exists and is verified
	isValidPartner, err := s.authUserClient.ValidatePartnerExists(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("validate partner: %w", err)
	}
	if !isValidPartner {
		return nil, fmt.Errorf("partner not found or not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("room not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	if !room.IsActive {
		return nil, fmt.Errorf("cannot create availability for inactive room")
	}

	// Check partner has access to the room at the specified time
	hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, partnerID, roomID, startTime)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("partner does not have allocation for this room")
		}
		return nil, fmt.Errorf("check partner room access: %w", err)
	}

	if !hasAccess.IsActiveAt(startTime) || !hasAccess.IsActiveAt(endTime) {
		return nil, fmt.Errorf("partner does not have access to room during specified time")
	}

	// Check for scheduling conflicts
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, partnerID, startTime, endTime, nil)
	if err != nil {
		return nil, fmt.Errorf("check availability conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("availability conflicts with existing schedule")
	}

	// Create domain entity with validation
	availability, err := domain.NewAvailability(partnerID, roomID, startTime, endTime, maxCapacity)
	if err != nil {
		return nil, fmt.Errorf("create availability entity: %w", err)
	}

	// Persist to repository
	if err := s.availabilityRepo.Create(ctx, availability); err != nil {
		return nil, fmt.Errorf("create availability: %w", err)
	}

	return availability, nil
}
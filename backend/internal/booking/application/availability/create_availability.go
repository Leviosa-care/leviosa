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

func (s *AvailabilityService) CreateAvailability(ctx context.Context, request *domain.CreateAvailabilityRequest) (*domain.Availability, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Verify roomEncx exists and is active
	roomEncx, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("get room by ID to verify existence: %w", err)
	}

	if !roomEncx.IsActive {
		return nil, fmt.Errorf("cannot create availability for inactive room")
	}

	// Check partner has access to the room at the specified time
	hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, partnerID, roomID, startTime)
	if err != nil {
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

	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, s.crypto, availability)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("availability", err)
	}

	// Persist to repository
	if err := s.availabilityRepo.Create(ctx, availabilityEncx); err != nil {
		return nil, fmt.Errorf("create availability: %w", err)
	}

	return availability, nil
}

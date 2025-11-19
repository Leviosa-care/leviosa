package availability

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AvailabilityService) CreateRecurringAvailability(ctx context.Context, request *domain.CreateRecurringAvailabilityRequest) (*domain.Availability, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Verify roomEncx exists and is active
	roomEncx, err := s.roomRepo.GetByID(ctx, request.RoomID)
	if err != nil {
		return nil, fmt.Errorf("get room by ID to verify room existence: %w", err)
	}

	if !roomEncx.IsActive {
		return nil, fmt.Errorf("cannot create availability for inactive room")
	}

	// Check partner has access to the room
	hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, request.UserID, request.RoomID, request.StartTime)
	if err != nil {
		return nil, fmt.Errorf("check partner room access: %w", err)
	}

	if !hasAccess.IsActiveAt(request.StartTime) {
		return nil, fmt.Errorf("partner does not have access to room during specified time")
	}

	// Check for scheduling conflicts
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, request.UserID, request.StartTime, request.EndTime, nil)
	if err != nil {
		return nil, fmt.Errorf("check availability conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("recurring availability conflicts with existing schedule")
	}

	// Create domain entity with validation
	availability, err := domain.NewRecurringAvailability(request.UserID, request.RoomID, request.StartTime, request.EndTime, request.MaxCapacity, request.Pattern)
	if err != nil {
		return nil, fmt.Errorf("create recurring availability entity: %w", err)
	}

	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, s.crypto, availability)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("availability", err)
	}

	// Persist to repository
	if err := s.availabilityRepo.Create(ctx, availabilityEncx); err != nil {
		return nil, fmt.Errorf("create recurring availability: %w", err)
	}

	return availability, nil
}

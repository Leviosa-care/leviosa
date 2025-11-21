package availability

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *AvailabilityService) CreateAvailability(ctx context.Context, request *domain.CreateAvailabilityRequest) (*domain.Availability, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Verify roomEncx exists and is active
	roomEncx, err := s.roomRepo.GetByID(ctx, request.RoomID)
	if err != nil {
		return nil, fmt.Errorf("get room by ID to verify existence: %w", err)
	}

	if !roomEncx.IsActive {
		return nil, fmt.Errorf("cannot create availability for inactive room")
	}

	userIDBytes, err := encx.SerializeValue(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("serialize user ID for hashing: %w", err)
	}
	userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

	// Check partner has access to the room at the specified time
	roomAllocationEncx, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, request.RoomID, request.StartTime)
	if err != nil {
		return nil, fmt.Errorf("check partner room access: %w", err)
	}

	hasAccess, err := domain.DecryptRoomAllocationEncx(ctx, s.crypto, roomAllocationEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("room allocation", err)
	}

	if !hasAccess.IsActiveAt(request.StartTime) || !hasAccess.IsActiveAt(request.EndTime) {
		return nil, fmt.Errorf("partner does not have access to room during specified time")
	}

	// Check for scheduling conflicts
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, request.UserID, request.StartTime, request.EndTime, nil)
	if err != nil {
		return nil, fmt.Errorf("check availability conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("availability conflicts with existing schedule")
	}

	// Create domain entity with validation
	availability, err := domain.NewAvailability(request.UserID, request.RoomID, request.StartTime, request.EndTime, request.MaxCapacity)
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

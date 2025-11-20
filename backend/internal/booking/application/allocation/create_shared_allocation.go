package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// CreateSharedAllocation creates a shared room allocation for a partner
func (s *RoomAllocationService) CreateSharedAllocation(ctx context.Context, request *domain.CreateSharedAllocationRequest) (*domain.RoomAllocation, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Validate partner exists and is verified
	partner, err := s.authUserClient.GetPartnerByUserID(ctx, request.UserID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewInvalidInputErr(errors.New("partner by user ID not found"))
		}
		return nil, err
	}
	if !partner.IsVerified {
		return nil, fmt.Errorf("partner is not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, request.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewInvalidInputErr(errors.New("room with ID " + request.RoomID.String() + " not found"))
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	if !room.IsActive {
		return nil, fmt.Errorf("cannot allocate inactive room")
	}

	// Check for existing allocation conflict
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, request.RoomID, request.UserID, domain.AllocationTypeShared, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("partner already has allocation for this room")
	}

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(request.RoomID, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("create shared allocation entity: %w", err)
	}

	// Persist to repository
	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("create shared allocation: %w", err)
	}

	return allocation, nil
}

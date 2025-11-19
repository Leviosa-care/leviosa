package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// CreateSharedAllocation creates a shared room allocation for a partner
func (s *RoomAllocationService) CreateSharedAllocation(ctx context.Context, roomID, partnerID uuid.UUID) (*domain.RoomAllocation, error) {
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
		return nil, fmt.Errorf("cannot allocate inactive room")
	}

	// Check for existing allocation conflict
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, roomID, partnerID, domain.AllocationTypeShared, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("partner already has allocation for this room")
	}

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(roomID, partnerID)
	if err != nil {
		return nil, fmt.Errorf("create shared allocation entity: %w", err)
	}

	// Persist to repository
	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("create shared allocation: %w", err)
	}

	return allocation, nil
}

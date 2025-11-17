package allocation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// GetAllocation retrieves a room allocation by ID
func (s *RoomAllocationService) GetAllocation(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error) {
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get allocation: %w", err)
	}

	return allocation, nil
}

// GetPartnerAllocations retrieves all allocations for a specific partner
func (s *RoomAllocationService) GetPartnerAllocations(ctx context.Context, partnerID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	allocations, err := s.allocationRepo.GetByUserID(ctx, partnerID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get partner allocations: %w", err)
	}

	return allocations, nil
}

// GetRoomAllocations retrieves all allocations for a specific room
func (s *RoomAllocationService) GetRoomAllocations(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	// Verify room exists
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	allocations, err := s.allocationRepo.GetByRoomID(ctx, roomID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get room allocations: %w", err)
	}

	return allocations, nil
}

// CheckPartnerRoomAccess checks if a partner has access to a room at a specific time
func (s *RoomAllocationService) CheckPartnerRoomAccess(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (bool, error) {
	// Get active allocation for partner and room at the specified time
	allocation, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, partnerID, roomID, at)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return false, nil // No allocation means no access
		}
		return false, fmt.Errorf("check partner room access: %w", err)
	}

	// Check if allocation is active at the specified time
	return allocation.IsActiveAt(at), nil
}
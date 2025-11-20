package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

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

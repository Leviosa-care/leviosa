package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetRoomAllocations retrieves all allocations for a specific room
func (s *RoomAllocationService) GetRoomAllocations(ctx context.Context, request *domain.GetRoomAllocationsRequest) ([]*domain.RoomAllocation, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Verify room exists
	_, err := s.roomRepo.GetByID(ctx, request.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewInvalidInputErr(errors.New("room with ID " + request.RoomID.String() + " not found"))
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	allocationsEncx, err := s.allocationRepo.GetByRoomID(ctx, request.RoomID, request.ActiveOnly)
	if err != nil {
		return nil, fmt.Errorf("get room allocations: %w", err)
	}

	// Decrypt all results
	allocations := make([]*domain.RoomAllocation, 0, len(allocationsEncx))
	for _, allocationEncx := range allocationsEncx {
		allocation, err := domain.DecryptRoomAllocationEncx(ctx, s.crypto, allocationEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt allocation: %w", err)
		}
		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

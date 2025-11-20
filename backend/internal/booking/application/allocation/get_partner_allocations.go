package allocation

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetPartnerAllocations retrieves all allocations for a specific partner
func (s *RoomAllocationService) GetPartnerAllocations(ctx context.Context, request *domain.GetPartnerAllocationsRequest) ([]*domain.RoomAllocation, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Compute hash for lookup
	userIDBytes, err := request.UserID.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("serialize user ID for hashing: %w", err)
	}
	userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

	allocationsEncx, err := s.allocationRepo.GetByUserIDHash(ctx, userIDHash, request.ActiveOnly)
	if err != nil {
		return nil, fmt.Errorf("get partner allocations: %w", err)
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

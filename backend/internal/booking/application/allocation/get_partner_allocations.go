package allocation

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// GetPartnerAllocations retrieves all allocations for a specific partner
func (s *RoomAllocationService) GetPartnerAllocations(ctx context.Context, partnerID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	allocations, err := s.allocationRepo.GetByUserID(ctx, partnerID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get partner allocations: %w", err)
	}

	return allocations, nil
}

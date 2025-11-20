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

	allocations, err := s.allocationRepo.GetByUserID(ctx, request.UserID, request.ActiveOnly)
	if err != nil {
		return nil, fmt.Errorf("get partner allocations: %w", err)
	}

	return allocations, nil
}

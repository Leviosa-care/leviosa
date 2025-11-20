package allocation

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"

	"github.com/google/uuid"
)

// GetAllocation retrieves a room allocation by ID
func (s *RoomAllocationService) GetAllocation(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error) {
	allocationEncx, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get allocation: %w", err)
	}

	// Decrypt before returning
	allocation, err := domain.DecryptRoomAllocationEncx(ctx, s.crypto, allocationEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt allocation: %w", err)
	}

	return allocation, nil
}

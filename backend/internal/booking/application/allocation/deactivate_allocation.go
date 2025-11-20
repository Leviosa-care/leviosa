package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeactivateAllocation deactivates a room allocation
func (s *RoomAllocationService) DeactivateAllocation(ctx context.Context, id uuid.UUID) error {
	// Get existing allocation
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewInvalidInputErr(errors.New("allocation by ID not found"))
		}
		return fmt.Errorf("get allocation for deactivation: %w", err)
	}

	// Deactivate
	allocation.Deactivate()

	// Persist changes
	if err := s.allocationRepo.Update(ctx, allocation); err != nil {
		return fmt.Errorf("deactivate allocation: %w", err)
	}

	return nil
}

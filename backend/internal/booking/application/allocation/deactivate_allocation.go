package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// DeactivateAllocation deactivates a room allocation
func (s *RoomAllocationService) DeactivateAllocation(ctx context.Context, id uuid.UUID) error {
	// Get existing allocation
	allocationEncx, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewInvalidInputErr(errors.New("allocation by ID not found"))
		}
		return fmt.Errorf("get allocation for deactivation: %w", err)
	}

	// Decrypt
	allocation, err := domain.DecryptRoomAllocationEncx(ctx, s.crypto, allocationEncx)
	if err != nil {
		return fmt.Errorf("decrypt allocation: %w", err)
	}

	// Deactivate
	allocation.Deactivate()

	// Re-encrypt before persisting
	allocationEncx, err = domain.ProcessRoomAllocationEncx(ctx, s.crypto, allocation)
	if err != nil {
		return fmt.Errorf("encrypt allocation: %w", err)
	}

	// Persist changes
	if err := s.allocationRepo.Update(ctx, allocationEncx); err != nil {
		return fmt.Errorf("deactivate allocation: %w", err)
	}

	return nil
}

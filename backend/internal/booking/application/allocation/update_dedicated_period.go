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

// UpdateDedicatedPeriod updates the time period for a dedicated allocation
func (s *RoomAllocationService) UpdateDedicatedPeriod(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error) {
	// Get existing allocation
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewInvalidInputErr(errors.New("allocation by ID not found"))
		}
		return nil, fmt.Errorf("get allocation for update: %w", err)
	}

	// Validate this is a dedicated allocation
	if allocation.AllocationType != domain.AllocationTypeDedicated {
		return nil, fmt.Errorf("can only update period for dedicated allocations")
	}

	// Validate dates
	if startDate == nil {
		return nil, fmt.Errorf("start date is required")
	}

	if endDate != nil && endDate.Before(*startDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Check for conflicts with other allocations (excluding this one)
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, allocation.RoomID, allocation.UserID, domain.AllocationTypeDedicated, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, errs.NewAlreadyExistsError(fmt.Errorf("conflicts with existing allocation"), "dedicated allocation")
	}

	// Update period with validation
	if err := allocation.UpdateDedicatedPeriod(startDate, endDate); err != nil {
		return nil, fmt.Errorf("update dedicated period: %w", err)
	}

	// Persist changes
	if err := s.allocationRepo.Update(ctx, allocation); err != nil {
		return nil, fmt.Errorf("update allocation period: %w", err)
	}

	return allocation, nil
}

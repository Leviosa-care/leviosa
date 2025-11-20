package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// UpdateDedicatedPeriod updates the time period for a dedicated allocation
// func (s *RoomAllocationService) UpdateDedicatedPeriod(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error) {
func (s *RoomAllocationService) UpdateDedicatedPeriod(ctx context.Context, request *domain.UpdateDedicatedAllocationRequest) (*domain.RoomAllocation, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing allocation
	allocation, err := s.allocationRepo.GetByID(ctx, request.ID)
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
	if request.StartDate == nil {
		return nil, fmt.Errorf("start date is required")
	}

	if request.EndDate != nil && request.EndDate.Before(*request.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Check for conflicts with other allocations (excluding this one)
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, allocation.RoomID, allocation.UserID, domain.AllocationTypeDedicated, request.StartDate, request.EndDate, &request.ID)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, errs.NewAlreadyExistsError(fmt.Errorf("conflicts with existing allocation"), "dedicated allocation")
	}

	// Update period with validation
	if err := allocation.UpdateDedicatedPeriod(request.StartDate, request.EndDate); err != nil {
		return nil, fmt.Errorf("update dedicated period: %w", err)
	}

	// Persist changes
	if err := s.allocationRepo.Update(ctx, allocation); err != nil {
		return nil, fmt.Errorf("update allocation period: %w", err)
	}

	return allocation, nil
}

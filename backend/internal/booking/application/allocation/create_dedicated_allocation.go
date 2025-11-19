package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// CreateDedicatedAllocation creates a dedicated room allocation with time bounds
func (s *RoomAllocationService) CreateDedicatedAllocation(ctx context.Context, request *domain.CreateDedicatedAllocationRequest) (*domain.RoomAllocation, error) {
	// Validate partner exists and is verified
	isValidPartner, err := s.authUserClient.ValidatePartnerExists(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("validate partner: %w", err)
	}
	if !isValidPartner {
		return nil, fmt.Errorf("partner not found or not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, request.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("room not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	if !room.IsActive {
		return nil, fmt.Errorf("cannot allocate inactive room")
	}

	// Validate dates
	if request.StartDate == nil {
		return nil, fmt.Errorf("start date is required for dedicated allocations")
	}

	if request.EndDate != nil && request.EndDate.Before(*request.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Check for existing allocation conflict
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, request.RoomID, request.UserID, domain.AllocationTypeDedicated, request.StartDate, request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("dedicated allocation conflicts with existing allocation")
	}

	// Create domain entity
	allocation, err := domain.NewDedicatedAllocation(request.RoomID, request.UserID, request.StartDate, request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("create dedicated allocation entity: %w", err)
	}

	// Persist to repository
	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("create dedicated allocation: %w", err)
	}

	return allocation, nil
}

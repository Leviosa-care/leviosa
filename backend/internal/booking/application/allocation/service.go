package allocation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

type RoomAllocationService struct {
	allocationRepo ports.RoomAllocationRepository
	roomRepo       ports.RoomRepository
	authUserClient ports.AuthUserClient
}

// New creates a new instance of the room allocation service
func New(allocationRepo ports.RoomAllocationRepository, roomRepo ports.RoomRepository, authUserClient ports.AuthUserClient) ports.RoomAllocationService {
	return &RoomAllocationService{
		allocationRepo: allocationRepo,
		roomRepo:       roomRepo,
		authUserClient: authUserClient,
	}
}

func (s *RoomAllocationService) CreateSharedAllocation(ctx context.Context, roomID, partnerID uuid.UUID) (*domain.RoomAllocation, error) {
	// Validate partner exists and is verified
	isValidPartner, err := s.authUserClient.ValidatePartnerExists(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("validate partner: %w", err)
	}
	if !isValidPartner {
		return nil, fmt.Errorf("partner not found or not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("room not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	if !room.IsActive {
		return nil, fmt.Errorf("cannot allocate inactive room")
	}

	// Check for existing allocation conflict
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, roomID, partnerID, domain.AllocationTypeShared, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("partner already has allocation for this room")
	}

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(roomID, partnerID)
	if err != nil {
		return nil, fmt.Errorf("create shared allocation entity: %w", err)
	}

	// Persist to repository
	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("create shared allocation: %w", err)
	}

	return allocation, nil
}

func (s *RoomAllocationService) CreateDedicatedAllocation(ctx context.Context, roomID, partnerID uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error) {
	// Validate partner exists and is verified
	isValidPartner, err := s.authUserClient.ValidatePartnerExists(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("validate partner: %w", err)
	}
	if !isValidPartner {
		return nil, fmt.Errorf("partner not found or not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, roomID)
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
	if startDate == nil {
		return nil, fmt.Errorf("start date is required for dedicated allocations")
	}

	if endDate != nil && endDate.Before(*startDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Check for existing allocation conflict
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, roomID, partnerID, domain.AllocationTypeDedicated, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("dedicated allocation conflicts with existing allocation")
	}

	// Create domain entity
	var allocation *domain.RoomAllocation
	if endDate != nil {
		allocation, err = domain.NewDedicatedAllocation(roomID, partnerID, *startDate, *endDate)
	} else {
		allocation, err = domain.NewDedicatedAllocation(roomID, partnerID, *startDate, time.Time{})
	}
	if err != nil {
		return nil, fmt.Errorf("create dedicated allocation entity: %w", err)
	}

	// Persist to repository
	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("create dedicated allocation: %w", err)
	}

	return allocation, nil
}

func (s *RoomAllocationService) GetAllocation(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error) {
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get allocation: %w", err)
	}

	return allocation, nil
}

func (s *RoomAllocationService) UpdateDedicatedPeriod(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error) {
	// Get existing allocation
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
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
	hasConflict, err := s.allocationRepo.CheckConflict(ctx, allocation.RoomID, allocation.PartnerID, domain.AllocationTypeDedicated, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("check allocation conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("updated period conflicts with existing allocation")
	}

	// Update period with validation
	var endTime time.Time
	if endDate != nil {
		endTime = *endDate
	}
	if err := allocation.UpdateDedicatedPeriod(*startDate, endTime); err != nil {
		return nil, fmt.Errorf("update dedicated period: %w", err)
	}

	// Persist changes
	if err := s.allocationRepo.Update(ctx, allocation); err != nil {
		return nil, fmt.Errorf("update allocation period: %w", err)
	}

	return allocation, nil
}

func (s *RoomAllocationService) DeactivateAllocation(ctx context.Context, id uuid.UUID) error {
	// Get existing allocation
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
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

func (s *RoomAllocationService) GetPartnerAllocations(ctx context.Context, partnerID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	allocations, err := s.allocationRepo.GetByPartnerID(ctx, partnerID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get partner allocations: %w", err)
	}

	return allocations, nil
}

func (s *RoomAllocationService) GetRoomAllocations(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	// Verify room exists
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	allocations, err := s.allocationRepo.GetByRoomID(ctx, roomID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get room allocations: %w", err)
	}

	return allocations, nil
}

func (s *RoomAllocationService) CheckPartnerRoomAccess(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (bool, error) {
	// Get active allocation for partner and room at the specified time
	allocation, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, partnerID, roomID, at)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return false, nil // No allocation means no access
		}
		return false, fmt.Errorf("check partner room access: %w", err)
	}

	// Check if allocation is active at the specified time
	return allocation.IsActiveAt(at), nil
}

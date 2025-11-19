package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// RoomAllocationService defines the interface for room allocation business logic
type RoomAllocationService interface {
	// CreateSharedAllocation creates a shared room allocation
	CreateSharedAllocation(ctx context.Context, request *domain.CreateSharedAllocationRequest) (*domain.RoomAllocation, error)

	// CreateDedicatedAllocation creates a dedicated room allocation with time bounds
	CreateDedicatedAllocation(ctx context.Context, request *domain.CreateDedicatedAllocationRequest) (*domain.RoomAllocation, error)

	// GetAllocation retrieves an allocation by ID
	GetAllocation(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error)

	// UpdateDedicatedPeriod updates the time period for a dedicated allocation
	UpdateDedicatedPeriod(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error)

	// DeactivateAllocation marks an allocation as inactive
	DeactivateAllocation(ctx context.Context, id uuid.UUID) error

	// GetPartnerAllocations retrieves all allocations for a partner
	GetPartnerAllocations(ctx context.Context, partnerID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error)

	// GetRoomAllocations retrieves all allocations for a room
	GetRoomAllocations(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error)

	// CheckPartnerRoomAccess verifies if a partner has access to a room at a specific time
	CheckPartnerRoomAccess(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (bool, error)
}

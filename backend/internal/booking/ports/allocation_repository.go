package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"

	"github.com/google/uuid"
)

// RoomAllocationRepository defines the interface for room allocation data persistence
type RoomAllocationRepository interface {
	// Create stores a new room allocation
	Create(ctx context.Context, allocation *domain.RoomAllocation) error

	// GetByID retrieves a room allocation by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error)

	// Update modifies an existing room allocation
	Update(ctx context.Context, allocation *domain.RoomAllocation) error

	// Delete removes a room allocation (soft delete by marking inactive)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves room allocations with optional filtering
	List(ctx context.Context, filter RoomAllocationFilter) ([]*domain.RoomAllocation, error)

	// GetByUserID retrieves all allocations for a specific partner
	GetByUserID(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error)

	// GetByRoomID retrieves all allocations for a specific room
	GetByRoomID(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error)

	// GetActiveAllocationForPartnerAndRoom checks if a partner has an active allocation for a room at a given time
	GetActiveAllocationForPartnerAndRoom(ctx context.Context, userID, roomID uuid.UUID, at time.Time) (*domain.RoomAllocation, error)

	// CheckConflict checks if a new allocation would conflict with existing ones
	// excludeID allows excluding a specific allocation from the conflict check (useful for updates)
	CheckConflict(ctx context.Context, roomID, userID uuid.UUID, allocationType domain.AllocationType, startDate, endDate *time.Time, excludeID *uuid.UUID) (bool, error)
}

// RoomAllocationFilter defines filtering options for room allocation queries
type RoomAllocationFilter struct {
	// Entity filters
	RoomID *uuid.UUID
	UserID *uuid.UUID

	// Allocation type filter
	AllocationType *domain.AllocationType

	// Active status filter
	IsActive *bool

	// Time-based filters
	ActiveAt     *time.Time // Check if allocation is active at this time
	OverlapsWith *TimeRange // Check if allocation overlaps with this time range

	// Pagination
	Limit  int
	Offset int

	// Sorting
	OrderBy        string // "created_at", "start_date", "end_date"
	OrderDirection string // "asc", "desc"
}

// TimeRange represents a time period for filtering
type TimeRange struct {
	Start time.Time
	End   time.Time
}

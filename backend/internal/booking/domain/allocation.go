package domain

import (
	"time"

	"github.com/google/uuid"
)

// AllocationType defines how a partner can access a room
type AllocationType string

const (
	AllocationTypeDedicated AllocationType = "dedicated" // Partner has exclusive access during specified period
	AllocationTypeShared    AllocationType = "shared"    // Partner shares room access with others
)

// RoomAllocation represents a partner's assignment to a room
type RoomAllocation struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"room_id"`
	UserID uuid.UUID `json:"user_id"`

	// Allocation configuration
	AllocationType AllocationType `json:"allocation_type"`

	// Time-based allocation (for dedicated allocations)
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`

	// Administrative fields
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewSharedAllocation creates a new shared room allocation
func NewSharedAllocation(roomID, userID uuid.UUID) (*RoomAllocation, error) {
	if roomID == uuid.Nil {
		return nil, ErrInvalidRoomID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidPartnerID
	}

	return &RoomAllocation{
		ID:             uuid.New(),
		RoomID:         roomID,
		UserID:         userID,
		AllocationType: AllocationTypeShared,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// NewDedicatedAllocation creates a new dedicated room allocation with time bounds
func NewDedicatedAllocation(roomID, userID uuid.UUID, startDate, endDate *time.Time) (*RoomAllocation, error) {
	if roomID == uuid.Nil {
		return nil, ErrInvalidRoomID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidPartnerID
	}
	if startDate == nil || startDate.IsZero() {
		return nil, ErrInvalidAllocationStartDate
	}
	if endDate != nil && endDate.Before(*startDate) {
		return nil, ErrInvalidAllocationEndDate
	}

	return &RoomAllocation{
		ID:             uuid.New(),
		RoomID:         roomID,
		UserID:         userID,
		AllocationType: AllocationTypeDedicated,
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// UpdateDedicatedPeriod updates the time period for a dedicated allocation
func (ra *RoomAllocation) UpdateDedicatedPeriod(startDate, endDate *time.Time) error {
	if ra.AllocationType != AllocationTypeDedicated {
		return ErrCannotUpdateSharedAllocationPeriod
	}
	if startDate == nil || startDate.IsZero() {
		return ErrInvalidAllocationStartDate
	}
	if endDate != nil && endDate.Before(*startDate) {
		return ErrInvalidAllocationEndDate
	}

	ra.StartDate = startDate
	ra.EndDate = endDate
	ra.UpdatedAt = time.Now()
	return nil
}

// IsActiveAt checks if the allocation is active at a given time
func (ra *RoomAllocation) IsActiveAt(t time.Time) bool {
	if !ra.IsActive {
		return false
	}

	// For shared allocations, they're always active if the allocation itself is active
	if ra.AllocationType == AllocationTypeShared {
		return true
	}

	// For dedicated allocations, check time bounds
	if ra.StartDate != nil && t.Before(*ra.StartDate) {
		return false
	}
	if ra.EndDate != nil && t.After(*ra.EndDate) {
		return false
	}

	return true
}

// Deactivate marks the allocation as inactive
func (ra *RoomAllocation) Deactivate() {
	ra.IsActive = false
	ra.UpdatedAt = time.Now()
}

// Activate marks the allocation as active
func (ra *RoomAllocation) Activate() {
	ra.IsActive = true
	ra.UpdatedAt = time.Now()
}

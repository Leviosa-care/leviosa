package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// AvailabilityRepository defines the interface for availability data persistence
type AvailabilityRepository interface {
	// Create stores a new availability
	Create(ctx context.Context, availability *domain.AvailabilityEncx) error

	// GetByID retrieves an availability by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AvailabilityEncx, error)

	// Update modifies an existing availability
	Update(ctx context.Context, availability *domain.AvailabilityEncx) error

	// Delete removes an availability
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves availabilities with optional filtering
	List(ctx context.Context, filter AvailabilityFilter) ([]*domain.AvailabilityEncx, error)

	// GetByuserID retrieves all availabilities for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID, filter AvailabilityFilter) ([]*domain.AvailabilityEncx, error)

	// GetByRoomID retrieves all availabilities for a specific room
	GetByRoomID(ctx context.Context, roomID uuid.UUID, filter AvailabilityFilter) ([]*domain.AvailabilityEncx, error)

	// GetAvailableSlots retrieves available slots within a time range
	GetAvailableSlots(ctx context.Context, filter AvailabilityFilter) ([]*domain.AvailabilityEncx, error)

	// CheckConflict checks if a new availability would conflict with existing ones for the same user
	CheckConflict(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error)

	// GetRecurringAvailabilities retrieves recurring availabilities that need to be expanded
	GetRecurringAvailabilities(ctx context.Context, until time.Time) ([]*domain.AvailabilityEncx, error)
}

// AvailabilityFilter defines filtering options for availability queries
type AvailabilityFilter struct {
	// Entity filters
	UserID *uuid.UUID
	RoomID *uuid.UUID

	// Status filter
	Status []domain.AvailabilityStatus

	// Time-based filters
	StartTime     *time.Time // Availabilities starting after this time
	EndTime       *time.Time // Availabilities ending before this time
	TimeRange     *TimeRange // Availabilities overlapping with this range
	AvailableOnly bool       // Only available slots

	// Service filters
	ServiceType *string
	MinPrice    *int // In cents
	MaxPrice    *int // In cents

	// Recurrence filter
	IsRecurring *bool

	// Pagination
	Limit  int
	Offset int

	// Sorting
	OrderBy        string // "start_time", "end_time", "created_at", "price_cents"
	OrderDirection string // "asc", "desc"
}

package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// RoomRepository defines the interface for room data persistence
type RoomRepository interface {
	// Create stores a new room
	Create(ctx context.Context, room *domain.RoomEncx) error

	// GetByID retrieves a room by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomEncx, error)

	// Update modifies an existing room
	Update(ctx context.Context, room *domain.RoomEncx) error

	// Delete removes a room (soft delete by marking inactive)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves rooms with optional filtering
	List(ctx context.Context, filter RoomFilter) ([]*domain.RoomEncx, error)

	// Count returns the total number of rooms matching the filter
	Count(ctx context.Context, filter RoomFilter) (int, error)

	// GetByBuildingID retrieves all rooms in a specific building
	GetByBuildingID(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.RoomEncx, error)
}

// RoomFilter defines filtering options for room queries
type RoomFilter struct {
	// Building filter
	BuildingID *uuid.UUID

	// Active status filter
	IsActive *bool

	// Capacity filter
	MinCapacity *int
	MaxCapacity *int

	// Searchable fields (plaintext - used by handler)
	Name       *string
	RoomNumber *string

	// Searchable field hashes (used by repository for database queries)
	NameHash       *string
	RoomNumberHash *string

	// Pagination
	Limit  int
	Offset int

	// Sorting
	OrderBy        string // "name", "created_at", "capacity"
	OrderDirection string // "asc", "desc"
}

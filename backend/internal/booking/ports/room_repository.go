package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// RoomRepository defines the interface for room data persistence
type RoomRepository interface {
	// Create stores a new room
	Create(ctx context.Context, room *domain.Room) error

	// GetByID retrieves a room by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error)

	// Update modifies an existing room
	Update(ctx context.Context, room *domain.Room) error

	// Delete removes a room (soft delete by marking inactive)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves rooms with optional filtering
	List(ctx context.Context, filter RoomFilter) ([]*domain.Room, error)

	// Count returns the total number of rooms matching the filter
	Count(ctx context.Context, filter RoomFilter) (int, error)

	// GetByBuildingID retrieves all rooms in a specific building
	GetByBuildingID(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.Room, error)
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

	// Rate filter (in cents)
	MinHourlyRate *int
	MaxHourlyRate *int

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
	OrderBy        string // "name", "created_at", "capacity", "hourly_rate_cents"
	OrderDirection string // "asc", "desc"
}

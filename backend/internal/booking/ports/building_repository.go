package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// BuildingRepository defines the interface for building data persistence
type BuildingRepository interface {
	// Create stores a new building
	Create(ctx context.Context, building *domain.BuildingEncx) error

	// GetByID retrieves a building by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BuildingEncx, error)

	// Update modifies an existing building
	Update(ctx context.Context, building *domain.BuildingEncx) error

	// Delete removes a building (soft delete by marking inactive)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves buildings with optional filtering
	List(ctx context.Context, filter BuildingFilter) ([]*domain.BuildingEncx, error)

	// Count returns the total number of buildings matching the filter
	Count(ctx context.Context, filter BuildingFilter) (int, error)

	// ExistsByNameOrAddress checks if a building with given name or address hash exists
	ExistsByNameOrAddress(ctx context.Context, nameHash, addressHash string) (bool, error)
}

// BuildingFilter defines filtering options for building queries
type BuildingFilter struct {
	// Active status filter
	IsActive *bool

	// Location filters (plaintext - used by handler)
	City    *string
	Country *string

	// Location filter hashes (used by repository for database queries)
	CityHash    *string
	CountryHash *string

	// Pagination
	Limit  int
	Offset int

	// Sorting
	OrderBy        string // "name", "created_at", "city"
	OrderDirection string // "asc", "desc"
}

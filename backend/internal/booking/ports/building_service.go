package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// BuildingService defines the interface for building business logic
type BuildingService interface {
	// CreateBuilding creates a new building with validation
	CreateBuilding(ctx context.Context, request *domain.CreateBuildingRequest) (*domain.BuildingResponse, error)

	// GetBuilding retrieves a building by ID
	GetBuildingByID(ctx context.Context, id uuid.UUID) (*domain.BuildingResponse, error)

	// ListBuildings retrieves buildings with filtering
	ListBuildings(ctx context.Context, filter BuildingFilter) ([]*domain.BuildingResponse, error)
	//
	// UpdateBuilding updates building details with validation
	UpdateBuilding(ctx context.Context, request *domain.UpdateBuildingRequest) (*domain.BuildingResponse, error)
	//
	// // SetBuildingContactInfo sets optional contact information
	// SetBuildingContactInfo(ctx context.Context, id uuid.UUID, description, phone, email string) (*domain.Building, error)
	//
	// // DeactivateBuilding marks a building as inactive
	// DeactivateBuilding(ctx context.Context, id uuid.UUID) error
	//
	// // ActivateBuilding marks a building as active
	// ActivateBuilding(ctx context.Context, id uuid.UUID) error
	//
	// // GetBuildingCount returns total count with filtering
	// GetBuildingCount(ctx context.Context, filter BuildingFilter) (int, error)
}

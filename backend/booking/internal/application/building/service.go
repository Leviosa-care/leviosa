package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

type BuildingService struct {
	buildingRepo ports.BuildingRepository
}

// New creates a new instance of the building service
func New(buildingRepo ports.BuildingRepository) ports.BuildingService {
	return &BuildingService{
		buildingRepo: buildingRepo,
	}
}

func (s *BuildingService) CreateBuilding(ctx context.Context, name, address, city, postalCode, country string) (*domain.Building, error) {
	// Create domain entity with validation
	building, err := domain.NewBuilding(name, address, city, postalCode, country)
	if err != nil {
		return nil, fmt.Errorf("create building entity: %w", err)
	}

	// Persist to repository
	if err := s.buildingRepo.Create(ctx, building); err != nil {
		return nil, fmt.Errorf("create building: %w", err)
	}

	return building, nil
}

func (s *BuildingService) GetBuilding(ctx context.Context, id uuid.UUID) (*domain.Building, error) {
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get building: %w", err)
	}

	return building, nil
}

func (s *BuildingService) UpdateBuilding(ctx context.Context, id uuid.UUID, name, address, city, postalCode, country string) (*domain.Building, error) {
	// Get existing building
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get building for update: %w", err)
	}

	// Update details with validation
	if err := building.UpdateDetails(name, address, city, postalCode, country); err != nil {
		return nil, fmt.Errorf("update building details: %w", err)
	}

	// Persist changes
	if err := s.buildingRepo.Update(ctx, building); err != nil {
		return nil, fmt.Errorf("update building: %w", err)
	}

	return building, nil
}

func (s *BuildingService) SetBuildingContactInfo(ctx context.Context, id uuid.UUID, description, phone, email string) (*domain.Building, error) {
	// Get existing building
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get building for contact update: %w", err)
	}

	// Set contact info
	building.SetContactInfo(description, phone, email)

	// Persist changes
	if err := s.buildingRepo.Update(ctx, building); err != nil {
		return nil, fmt.Errorf("update building contact info: %w", err)
	}

	return building, nil
}

func (s *BuildingService) DeactivateBuilding(ctx context.Context, id uuid.UUID) error {
	// Get existing building
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get building for deactivation: %w", err)
	}

	// Deactivate
	building.Deactivate()

	// Persist changes
	if err := s.buildingRepo.Update(ctx, building); err != nil {
		return fmt.Errorf("deactivate building: %w", err)
	}

	return nil
}

func (s *BuildingService) ActivateBuilding(ctx context.Context, id uuid.UUID) error {
	// Get existing building
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get building for activation: %w", err)
	}

	// Activate
	building.Activate()

	// Persist changes
	if err := s.buildingRepo.Update(ctx, building); err != nil {
		return fmt.Errorf("activate building: %w", err)
	}

	return nil
}

func (s *BuildingService) ListBuildings(ctx context.Context, filter ports.BuildingFilter) ([]*domain.Building, error) {
	buildings, err := s.buildingRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list buildings: %w", err)
	}

	return buildings, nil
}

func (s *BuildingService) GetBuildingCount(ctx context.Context, filter ports.BuildingFilter) (int, error) {
	count, err := s.buildingRepo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count buildings: %w", err)
	}

	return count, nil
}
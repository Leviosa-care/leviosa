package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

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


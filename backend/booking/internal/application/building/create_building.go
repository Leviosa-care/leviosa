package building

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
)

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


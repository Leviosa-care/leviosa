package building

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
)

func (s *BuildingService) ListBuildings(ctx context.Context, filter ports.BuildingFilter) ([]*domain.Building, error) {
	buildings, err := s.buildingRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list buildings: %w", err)
	}

	return buildings, nil
}
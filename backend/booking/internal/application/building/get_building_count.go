package building

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/ports"
)

func (s *BuildingService) GetBuildingCount(ctx context.Context, filter ports.BuildingFilter) (int, error) {
	count, err := s.buildingRepo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count buildings: %w", err)
	}

	return count, nil
}
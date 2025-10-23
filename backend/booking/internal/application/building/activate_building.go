package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

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


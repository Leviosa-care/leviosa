package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

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


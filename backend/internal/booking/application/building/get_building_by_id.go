package building

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *BuildingService) GetBuildingByID(ctx context.Context, id uuid.UUID) (*domain.BuildingResponse, error) {
	buildingEncx, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get building by ID: %w", err)
	}

	building, err := domain.DecryptBuildingEncx(ctx, s.crypto, buildingEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("building", err)
	}

	return building.ToResponse(), nil
}

package building

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *BuildingService) UpdateBuilding(ctx context.Context, request *domain.UpdateBuildingRequest) (*domain.BuildingResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing building
	buildingEncx, err := s.buildingRepo.GetByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("get building by ID: %w", err)
	}

	building, err := domain.DecryptBuildingEncx(ctx, s.crypto, buildingEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("building", err)
	}

	if request.Name != nil {
		building.Name = *request.Name
	}

	if request.Address != nil {
		building.Address = *request.Address
	}

	if request.City != nil {
		building.City = *request.City
	}

	if request.PostalCode != nil {
		building.PostalCode = *request.PostalCode
	}

	if request.Country != nil {
		building.Country = *request.Country
	}

	if request.Description != nil {
		building.Description = *request.Description
	}

	if request.Phone != nil {
		building.Phone = *request.Phone
	}

	if request.Email != nil {
		building.Email = *request.Email
	}

	if request.IsActive != nil {
		building.IsActive = *request.IsActive
	}

	buildingEncx, err = domain.ProcessBuildingEncx(ctx, s.crypto, building)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("building", err)
	}

	// Persist changes
	if err := s.buildingRepo.Update(ctx, buildingEncx); err != nil {
		return nil, fmt.Errorf("update building: %w", err)
	}

	return building.ToResponse(), nil
}

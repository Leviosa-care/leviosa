package building

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// func (s *BuildingService) CreateBuilding(ctx context.Context, request *domain.CreateBuildingRequest) (*domain.Building, error) {
func (s *BuildingService) CreateBuilding(ctx context.Context, request *domain.CreateBuildingRequest) (*domain.BuildingResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	now := time.Now()

	building := &domain.Building{
		ID:         uuid.New(),
		Name:       request.Name,
		Address:    request.Address,
		City:       request.City,
		PostalCode: request.PostalCode,
		Country:    request.Country,
		IsActive:   request.IsActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// TODO: maybe check if some similar building exists ?

	buildingEncx, err := domain.ProcessBuildingEncx(ctx, s.crypto, building)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("building", err)
	}

	// Persist to repository
	if err := s.buildingRepo.Create(ctx, buildingEncx); err != nil {
		// TODO: make a switch to handle errors here return nil, fmt.Errorf("create building: %w", err)
	}

	return &domain.BuildingResponse{
		ID:          building.ID,
		Name:        building.Name,
		Address:     building.Address,
		City:        building.City,
		PostalCode:  building.PostalCode,
		Country:     building.PostalCode,
		Description: building.Description,
		Phone:       building.Phone,
		Email:       building.Email,
	}, nil
}

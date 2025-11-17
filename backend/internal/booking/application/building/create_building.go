package building

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"

	"github.com/google/uuid"
)

func (s *BuildingService) CreateBuilding(ctx context.Context, request *domain.CreateBuildingRequest) (*domain.BuildingResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Generate hashes for duplicate checking
	nameBytes, err := encx.SerializeValue(request.Name)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid name value: %v", err))
	}
	nameHash := s.crypto.HashBasic(ctx, nameBytes)

	addressBytes, err := encx.SerializeValue(request.Address)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid address value: %v", err))
	}
	addressHash := s.crypto.HashBasic(ctx, addressBytes)

	// Check for duplicate name or address
	exists, err := s.buildingRepo.ExistsByNameOrAddress(ctx, nameHash, addressHash)
	if err != nil {
		return nil, fmt.Errorf("check duplicate building: %w", err)
	}
	if exists {
		return nil, errs.NewAlreadyExistsError(errors.New("building with this name or address already exists"), "")
	}

	now := time.Now()

	building := &domain.Building{
		ID:          uuid.New(),
		Name:        request.Name,
		Address:     request.Address,
		City:        request.City,
		PostalCode:  request.PostalCode,
		Country:     request.Country,
		Description: request.Description,
		Phone:       request.Phone,
		Email:       request.Email,
		IsActive:    request.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	buildingEncx, err := domain.ProcessBuildingEncx(ctx, s.crypto, building)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("building", err)
	}

	// Persist to repository
	if err := s.buildingRepo.Create(ctx, buildingEncx); err != nil {
		return nil, fmt.Errorf("create building: %w", err)
	}

	return &domain.BuildingResponse{
		ID:          building.ID,
		Name:        building.Name,
		Address:     building.Address,
		City:        building.City,
		PostalCode:  building.PostalCode,
		Country:     building.Country,
		Description: building.Description,
		Phone:       building.Phone,
		Email:       building.Email,
		IsActive:    building.IsActive,
	}, nil
}

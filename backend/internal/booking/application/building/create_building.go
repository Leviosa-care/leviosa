package building

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *BuildingService) CreateBuilding(ctx context.Context, request *domain.CreateBuildingRequest) (*domain.BuildingResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
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
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("building data: %v", err))
		case errors.Is(err, errs.ErrUniqueViolation):
			return nil, errs.NewAlreadyExistsError(err, "building with this name or address")
		case errors.Is(err, errs.ErrNotNullViolation):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("missing required data for building: %v", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid foreign key for building: %v", err))
		case errors.Is(err, errs.ErrCheckViolation):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("building data failed check constraint: %v", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for building: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error for building: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during building creation: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during building creation: %w", err))
		}
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

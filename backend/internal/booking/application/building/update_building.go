package building

import (
	"context"
	"errors"
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "building")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during building retrieval: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("query failed during building retrieval: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during building retrieval: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during building retrieval: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled error during building retrieval: %w", err))
		}
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "building")
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
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during building update: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during building update: %w", err))
		}
	}

	return building.ToResponse(), nil
}

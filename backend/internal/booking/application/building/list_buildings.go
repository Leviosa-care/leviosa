package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *BuildingService) ListBuildings(ctx context.Context, filter ports.BuildingFilter) ([]*domain.BuildingResponse, error) {
	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.City != nil {
		cityBytes, err := encx.SerializeValue(*filter.City)
		if err != nil {
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid city value: %v", err))
		}
		cityHash := s.crypto.HashBasic(ctx, cityBytes)
		repoFilter.CityHash = &cityHash
	}

	if filter.Country != nil {
		countryBytes, err := encx.SerializeValue(*filter.Country)
		if err != nil {
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid country value: %v", err))
		}
		countryHash := s.crypto.HashBasic(ctx, countryBytes)
		repoFilter.CountryHash = &countryHash
	}

	buildingsEncx, err := s.buildingRepo.List(ctx, repoFilter)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid filter parameters: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during building list retrieval: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for building list: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during building list retrieval: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during building list retrieval: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during building list retrieval: %w", err))
		}
	}

	var buildings []*domain.BuildingResponse
	for _, buildingEncx := range buildingsEncx {
		building, err := domain.DecryptBuildingEncx(ctx, s.crypto, buildingEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("building", err)
		}
		buildings = append(buildings, building.ToResponse())
	}
	return buildings, nil
}

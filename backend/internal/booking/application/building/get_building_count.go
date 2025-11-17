package building

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/hengadev/encx"
)

func (s *BuildingService) GetBuildingCount(ctx context.Context, filter ports.BuildingFilter) (int, error) {
	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.City != nil {
		cityBytes, err := encx.SerializeValue(*filter.City)
		if err != nil {
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid city value: %v", err))
		}
		cityHash := s.crypto.HashBasic(ctx, cityBytes)
		repoFilter.CityHash = &cityHash
	}

	if filter.Country != nil {
		countryBytes, err := encx.SerializeValue(*filter.Country)
		if err != nil {
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid country value: %v", err))
		}
		countryHash := s.crypto.HashBasic(ctx, countryBytes)
		repoFilter.CountryHash = &countryHash
	}

	count, err := s.buildingRepo.Count(ctx, repoFilter)
	if err != nil {
		return 0, fmt.Errorf("count buildings: %w", err)
	}

	return count, nil
}

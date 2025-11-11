package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *BuildingService) GetBuildingCount(ctx context.Context, filter ports.BuildingFilter) (int, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("get building count: %w", err)
	}

	logger.InfoContext(ctx, "Service: Getting building count",
		"operation", "get_building_count",
		"filter", filter)

	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.City != nil {
		cityBytes, err := encx.SerializeValue(*filter.City)
		if err != nil {
			logger.WarnContext(ctx, "Service: Invalid city value",
				"error", err,
				"operation", "get_building_count")
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid city value: %v", err))
		}
		cityHash := s.crypto.HashBasic(ctx, cityBytes)
		repoFilter.CityHash = &cityHash
	}

	if filter.Country != nil {
		countryBytes, err := encx.SerializeValue(*filter.Country)
		if err != nil {
			logger.WarnContext(ctx, "Service: Invalid country value",
				"error", err,
				"operation", "get_building_count")
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid country value: %v", err))
		}
		countryHash := s.crypto.HashBasic(ctx, countryBytes)
		repoFilter.CountryHash = &countryHash
	}

	count, err := s.buildingRepo.Count(ctx, repoFilter)
	if err != nil {
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			errorContext = "database connection failure"
		case errors.Is(err, errs.ErrTooManyConnections):
			errorContext = "too many database connections"
		case errors.Is(err, errs.ErrResourceExhausted):
			errorContext = "database resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled):
			errorContext = "query cancelled"
		case errors.Is(err, errs.ErrTransactionFailure):
			errorContext = "transaction failure"
		case errors.Is(err, errs.ErrDeadlock):
			errorContext = "database deadlock"
		case errors.Is(err, errs.ErrInvalidInput):
			errorContext = "invalid filter parameters"
		default:
			errorContext = "unexpected repository error"
		}

		logger.ErrorContext(ctx, "Service: Failed to count buildings",
			"error", err,
			"filter", filter,
			"operation", "get_building_count",
			"context", errorContext)

		return 0, fmt.Errorf("count buildings: %w", err)
	}

	logger.InfoContext(ctx, "Service: Building count retrieved successfully",
		"count", count,
		"filter", filter,
		"operation", "get_building_count")

	return count, nil
}

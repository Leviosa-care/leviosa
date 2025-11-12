package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *RoomService) GetRoomCount(ctx context.Context, filter ports.RoomFilter) (int, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("get room count: %w", err)
	}

	logger.InfoContext(ctx, "Service: Getting room count",
		"operation", "get_room_count",
		"filter", filter)

	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.Name != nil {
		nameBytes, err := encx.SerializeValue(*filter.Name)
		if err != nil {
			logger.WarnContext(ctx, "Service: Invalid name value",
				"error", err,
				"operation", "get_room_count")
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid name value: %v", err))
		}
		nameHash := s.crypto.HashBasic(ctx, nameBytes)
		repoFilter.NameHash = &nameHash
	}

	if filter.RoomNumber != nil {
		roomNumberBytes, err := encx.SerializeValue(*filter.RoomNumber)
		if err != nil {
			logger.WarnContext(ctx, "Service: Invalid room number value",
				"error", err,
				"operation", "get_room_count")
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid room number value: %v", err))
		}
		roomNumberHash := s.crypto.HashBasic(ctx, roomNumberBytes)
		repoFilter.RoomNumberHash = &roomNumberHash
	}

	count, err := s.roomRepo.Count(ctx, repoFilter)
	if err != nil {
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			errorContext = "invalid filter parameters"
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
		case errors.Is(err, errs.ErrDBQuery):
			errorContext = "repository query failed"
		case errors.Is(err, errs.ErrDatabase):
			errorContext = "database error"
		case errors.Is(err, errs.ErrContext):
			errorContext = "context error"
		default:
			errorContext = "unexpected repository error"
		}

		logger.ErrorContext(ctx, "Service: Failed to count rooms",
			"error", err,
			"filter", filter,
			"operation", "get_room_count",
			"context", errorContext)

		return 0, fmt.Errorf("count rooms: %w", err)
	}

	logger.InfoContext(ctx, "Service: Room count retrieved successfully",
		"count", count,
		"filter", filter,
		"operation", "get_room_count")

	return count, nil
}
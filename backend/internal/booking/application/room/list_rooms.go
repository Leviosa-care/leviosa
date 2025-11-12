package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *RoomService) ListRooms(ctx context.Context, filter ports.RoomFilter) ([]*domain.Room, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %w", err)
	}

	logger.InfoContext(ctx, "Service: Listing rooms",
		"operation", "list_rooms",
		"filter", filter)

	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.Name != nil {
		nameBytes, err := encx.SerializeValue(*filter.Name)
		if err != nil {
			logger.WarnContext(ctx, "Service: Invalid name value",
				"error", err,
				"operation", "list_rooms")
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid name value: %v", err))
		}
		nameHash := s.crypto.HashBasic(ctx, nameBytes)
		repoFilter.NameHash = &nameHash
	}

	if filter.RoomNumber != nil {
		roomNumberBytes, err := encx.SerializeValue(*filter.RoomNumber)
		if err != nil {
			logger.WarnContext(ctx, "Service: Invalid room number value",
				"error", err,
				"operation", "list_rooms")
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid room number value: %v", err))
		}
		roomNumberHash := s.crypto.HashBasic(ctx, roomNumberBytes)
		repoFilter.RoomNumberHash = &roomNumberHash
	}

	roomsEncx, err := s.roomRepo.List(ctx, repoFilter)
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

		logger.ErrorContext(ctx, "Service: Failed to list rooms",
			"error", err,
			"filter", filter,
			"operation", "list_rooms",
			"context", errorContext)

		return nil, fmt.Errorf("list rooms: %w", err)
	}

	var rooms []*domain.Room
	for _, roomEncx := range roomsEncx {
		room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
		if err != nil {
			logger.ErrorContext(ctx, "Service: Failed to decrypt room",
				"error", err,
				"room_id", roomEncx.ID,
				"operation", "list_rooms")
			return nil, errs.NewNotDecryptedErr("room", err)
		}
		rooms = append(rooms, room)
	}

	logger.InfoContext(ctx, "Service: Rooms listed successfully",
		"count", len(rooms),
		"operation", "list_rooms")

	return rooms, nil
}
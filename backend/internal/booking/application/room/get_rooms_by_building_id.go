package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *RoomService) GetRoomsByBuilding(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.RoomResponse, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get rooms by building: %w", err)
	}

	logger.InfoContext(ctx, "Service: Getting rooms by building",
		"building_id", buildingID,
		"active_only", activeOnly,
		"operation", "get_rooms_by_building")

	// Verify building exists
	buildingEncx, err := s.buildingRepo.GetByID(ctx, buildingID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.ErrRepositoryNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid building ID: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during building verification: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for building verification: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during building verification: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during building verification: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during building verification: %w", err))
		}
	}

	// Decrypt building for logging (optional, just to verify it's decryptable)
	_, err = domain.DecryptBuildingEncx(ctx, s.crypto, buildingEncx)
	if err != nil {
		logger.WarnContext(ctx, "Service: Failed to decrypt building during verification",
			"error", err,
			"building_id", buildingID,
			"operation", "get_rooms_by_building")
	}

	roomsEncx, err := s.roomRepo.GetByBuildingID(ctx, buildingID, activeOnly)
	if err != nil {
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			errorContext = "invalid input parameters"
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

		logger.ErrorContext(ctx, "Service: Failed to get rooms by building",
			"error", err,
			"building_id", buildingID,
			"active_only", activeOnly,
			"operation", "get_rooms_by_building",
			"context", errorContext)

		return nil, fmt.Errorf("get rooms by building: %w", err)
	}

	var rooms []*domain.RoomResponse
	for _, roomEncx := range roomsEncx {
		room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
		if err != nil {
			logger.ErrorContext(ctx, "Service: Failed to decrypt room",
				"error", err,
				"room_id", roomEncx.ID,
				"building_id", buildingID,
				"operation", "get_rooms_by_building")
			return nil, errs.NewNotDecryptedErr("room", err)
		}
		rooms = append(rooms, room.ToResponse())
	}

	logger.InfoContext(ctx, "Service: Rooms retrieved by building successfully",
		"count", len(rooms),
		"building_id", buildingID,
		"active_only", activeOnly,
		"operation", "get_rooms_by_building")

	return rooms, nil
}


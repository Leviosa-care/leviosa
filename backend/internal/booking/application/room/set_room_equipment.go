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

func (s *RoomService) SetRoomEquipment(ctx context.Context, id uuid.UUID, equipment []string) (*domain.RoomResponse, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("set room equipment: %w", err)
	}

	logger.InfoContext(ctx, "Service: Setting room equipment",
		"room_id", id,
		"equipment_count", len(equipment),
		"operation", "set_room_equipment")

	// Validate equipment input
	if len(equipment) == 0 {
		logger.WarnContext(ctx, "Service: Empty equipment array provided",
			"room_id", id,
			"operation", "set_room_equipment")
		// This is not necessarily an error - empty equipment is valid
	} else {
		for i, item := range equipment {
			if item == "" {
				logger.WarnContext(ctx, "Service: Empty equipment item found",
					"room_id", id,
					"item_index", i,
					"operation", "set_room_equipment")
			}
		}
	}

	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Room not found for equipment update",
				"room_id", id,
				"operation", "set_room_equipment")
			return nil, errs.NewNotFoundErr(err, "room")
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid room ID for equipment update: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during room retrieval for equipment update: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room equipment update: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during room retrieval for equipment update: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during room retrieval for equipment update: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room retrieval for equipment update: %w", err))
		}
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to decrypt room for equipment update",
			"error", err,
			"room_id", id,
			"operation", "set_room_equipment")
		return nil, errs.NewNotDecryptedErr("room", err)
	}

	// Set equipment
	room.SetEquipment(equipment)

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to encrypt room for equipment update",
			"error", err,
			"room_id", id,
			"operation", "set_room_equipment")
		return nil, errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Room not found during equipment update",
				"room_id", id,
				"operation", "set_room_equipment")
			return nil, errs.NewNotFoundErr(err, "room")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during room equipment update: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room equipment update: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during room equipment update: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during room equipment update: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room equipment update: %w", err))
		}
	}

	logger.InfoContext(ctx, "Service: Room equipment updated successfully",
		"room_id", id,
		"equipment_count", len(equipment),
		"operation", "set_room_equipment")

	return room.ToResponse(), nil
}


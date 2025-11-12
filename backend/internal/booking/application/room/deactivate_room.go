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

func (s *RoomService) DeactivateRoom(ctx context.Context, id uuid.UUID) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("deactivate room: %w", err)
	}

	logger.InfoContext(ctx, "Service: Deactivating room",
		"room_id", id,
		"operation", "deactivate_room")

	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Room not found for deactivation",
				"room_id", id,
				"operation", "deactivate_room")
			return errs.ErrRepositoryNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid room ID for deactivation: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error during room retrieval for deactivation: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room deactivation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database error during room retrieval for deactivation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during room retrieval for deactivation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room retrieval for deactivation: %w", err))
		}
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to decrypt room for deactivation",
			"error", err,
			"room_id", id,
			"operation", "deactivate_room")
		return errs.NewNotDecryptedErr("room", err)
	}

	// Check if room is already inactive
	if !room.IsActive {
		logger.InfoContext(ctx, "Service: Room is already inactive",
			"room_id", id,
			"operation", "deactivate_room")
		// This is not an error - room is already in desired state
		return nil
	}

	// Check for active allocations or bookings before deactivating
	// Note: This would typically involve checking for future bookings, room allocations, etc.
	// For now, we'll proceed with deactivation, but in a real implementation you might want to check:
	// - Active allocations
	// - Future bookings
	// - Whether the room can be safely deactivated

	logger.InfoContext(ctx, "Service: Room will be deactivated",
		"room_id", id,
		"room_name", room.Name,
		"room_number", room.RoomNumber,
		"operation", "deactivate_room")

	// Deactivate room
	room.Deactivate()

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to encrypt room for deactivation",
			"error", err,
			"room_id", id,
			"operation", "deactivate_room")
		return errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Room not found during deactivation",
				"room_id", id,
				"operation", "deactivate_room")
			return errs.ErrRepositoryNotFound
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error during room deactivation: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room deactivation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database error during room deactivation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during room deactivation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room deactivation: %w", err))
		}
	}

	logger.InfoContext(ctx, "Service: Room deactivated successfully",
		"room_id", id,
		"room_name", room.Name,
		"room_number", room.RoomNumber,
		"building_id", room.BuildingID,
		"operation", "deactivate_room")

	return nil
}
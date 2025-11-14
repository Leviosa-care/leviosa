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

func (s *RoomService) ActivateRoom(ctx context.Context, id uuid.UUID) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("activate room: %w", err)
	}

	logger.InfoContext(ctx, "Service: Activating room",
		"room_id", id,
		"operation", "activate_room")

	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Room not found for activation",
				"room_id", id,
				"operation", "activate_room")
			return errs.NewNotFoundErr(err, "room")
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid room ID for activation: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error during room retrieval for activation: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room activation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database error during room retrieval for activation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during room retrieval for activation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room retrieval for activation: %w", err))
		}
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to decrypt room for activation",
			"error", err,
			"room_id", id,
			"operation", "activate_room")
		return errs.NewNotDecryptedErr("room", err)
	}

	// Check if room is already active
	if room.IsActive {
		logger.InfoContext(ctx, "Service: Room is already active",
			"room_id", id,
			"operation", "activate_room")
		// This is not an error - room is already in desired state
		return nil
	}

	// Verify building is still active
	buildingEncx, err := s.buildingRepo.GetByID(ctx, room.BuildingID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Building not found for room activation",
				"room_id", id,
				"building_id", room.BuildingID,
				"operation", "activate_room")
			return errs.NewInvalidValueErr(fmt.Sprintf("cannot activate room in non-existent building %s", room.BuildingID))
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid building ID for room activation: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error during building verification for room activation: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for building verification for room activation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database error during building verification for room activation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during building verification for room activation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during building verification for room activation: %w", err))
		}
	}

	// Decrypt building
	building, err := domain.DecryptBuildingEncx(ctx, s.crypto, buildingEncx)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to decrypt building for room activation",
			"error", err,
			"room_id", id,
			"building_id", room.BuildingID,
			"operation", "activate_room")
		return errs.NewNotDecryptedErr("building", err)
	}

	if !building.IsActive {
		logger.WarnContext(ctx, "Service: Cannot activate room in inactive building",
			"room_id", id,
			"building_id", room.BuildingID,
			"building_name", building.Name,
			"operation", "activate_room")
		return errs.NewInvalidValueErr(fmt.Sprintf("cannot activate room %s in inactive building %s", id, room.BuildingID))
	}

	// Activate room
	room.Activate()

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		logger.ErrorContext(ctx, "Service: Failed to encrypt room for activation",
			"error", err,
			"room_id", id,
			"operation", "activate_room")
		return errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			logger.WarnContext(ctx, "Service: Room not found during activation",
				"room_id", id,
				"operation", "activate_room")
			return errs.NewNotFoundErr(err, "room")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error during room activation: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room activation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database error during room activation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during room activation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room activation: %w", err))
		}
	}

	logger.InfoContext(ctx, "Service: Room activated successfully",
		"room_id", id,
		"building_id", room.BuildingID,
		"building_name", building.Name,
		"operation", "activate_room")

	return nil
}
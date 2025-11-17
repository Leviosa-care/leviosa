package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *RoomService) ActivateRoom(ctx context.Context, id uuid.UUID) error {
	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get room by ID for activation: %w", err)
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("room", err)
	}

	// Verify building is still active
	buildingEncx, err := s.buildingRepo.GetByID(ctx, room.BuildingID)
	if err != nil {
		return fmt.Errorf("get building by ID for room activation: %w", err)
	}

	// Decrypt building
	building, err := domain.DecryptBuildingEncx(ctx, s.crypto, buildingEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("building", err)
	}

	if !building.IsActive {
		return fmt.Errorf("cannot activate room %s in inactive building %s", id, room.BuildingID)
	}

	// Activate room
	room.Activate()

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		return errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		return fmt.Errorf("update room by setting it active: %w", err)
	}

	return nil
}


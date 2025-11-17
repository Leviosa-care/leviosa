package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *RoomService) SetRoomEquipment(ctx context.Context, id uuid.UUID, equipment []string) (*domain.RoomResponse, error) {

	// Validate equipment input
	if len(equipment) == 0 {
		// This is not necessarily an error - empty equipment is valid
	} else {
		for i, item := range equipment {
			if item == "" {
				_ = i
			}
		}
	}

	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get room by ID for equipment update: %w", err)
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("room", err)
	}

	// Set equipment
	room.SetEquipment(equipment)

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		return nil, fmt.Errorf("set room equipment: %w", err)
	}

	return room.ToResponse(), nil
}

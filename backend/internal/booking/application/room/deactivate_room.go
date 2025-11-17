package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *RoomService) DeactivateRoom(ctx context.Context, id uuid.UUID) error {
	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get room by ID for deactivation: %w", err)
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("room", err)
	}

	// Check for active allocations or bookings before deactivating
	// Note: This would typically involve checking for future bookings, room allocations, etc.
	// For now, we'll proceed with deactivation, but in a real implementation you might want to check:
	// - Active allocations
	// - Future bookings
	// - Whether the room can be safely deactivated

	// Deactivate room
	room.Deactivate()

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		return errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		return fmt.Errorf("update room to set it as unactive: %w", err)
	}

	return nil
}


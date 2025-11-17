package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *RoomService) GetRoom(ctx context.Context, id uuid.UUID) (*domain.RoomResponse, error) {
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get room: %w", err)
	}

	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("room", err)
	}

	return room.ToResponse(), nil
}

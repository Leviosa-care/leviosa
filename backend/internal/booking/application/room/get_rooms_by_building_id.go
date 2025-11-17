package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *RoomService) GetRoomsByBuilding(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.RoomResponse, error) {
	roomsEncx, err := s.roomRepo.GetByBuildingID(ctx, buildingID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get rooms by building ID: %w", err)
	}

	rooms := []*domain.RoomResponse{}
	for _, roomEncx := range roomsEncx {
		room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("room", err)
		}
		rooms = append(rooms, room.ToResponse())
	}

	return rooms, nil
}

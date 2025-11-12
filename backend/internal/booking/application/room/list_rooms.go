package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *RoomService) ListRooms(ctx context.Context, filter ports.RoomFilter) ([]*domain.RoomResponse, error) {
	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.Name != nil {
		nameBytes, err := encx.SerializeValue(*filter.Name)
		if err != nil {
			println("HERE I HAVE AN ERROR IN THE FILTER NAME")
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid name value: %v", err))
		}
		nameHash := s.crypto.HashBasic(ctx, nameBytes)
		repoFilter.NameHash = &nameHash
	}

	if filter.RoomNumber != nil {
		roomNumberBytes, err := encx.SerializeValue(*filter.RoomNumber)
		if err != nil {
			println("HERE I HAVE AN ERROR")
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid room number value: %v", err))
		}
		roomNumberHash := s.crypto.HashBasic(ctx, roomNumberBytes)
		repoFilter.RoomNumberHash = &roomNumberHash
	}

	roomsEncx, err := s.roomRepo.List(ctx, repoFilter)
	if err != nil {
		println("HERE I HAVE AN ERROR IN THE LIST REPO FUNCTION")
		return nil, fmt.Errorf("list rooms: %w", err)
	}

	var rooms []*domain.RoomResponse
	for _, roomEncx := range roomsEncx {
		room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
		if err != nil {
			println("HERE I HAVE AN ERROR IN THE DECRYPT")
			return nil, errs.NewNotDecryptedErr("room", err)
		}
		rooms = append(rooms, room.ToResponse())
	}

	return rooms, nil
}

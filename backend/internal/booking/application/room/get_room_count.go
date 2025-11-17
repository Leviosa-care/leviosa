package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"
)

func (s *RoomService) GetRoomCount(ctx context.Context, filter ports.RoomFilter) (int, error) {
	// Create modified filter with hashed values for searchable fields
	repoFilter := filter

	if filter.Name != nil {
		nameBytes, err := encx.SerializeValue(*filter.Name)
		if err != nil {
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid name value: %v", err))
		}
		nameHash := s.crypto.HashBasic(ctx, nameBytes)
		repoFilter.NameHash = &nameHash
	}

	if filter.RoomNumber != nil {
		roomNumberBytes, err := encx.SerializeValue(*filter.RoomNumber)
		if err != nil {
			return 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid room number value: %v", err))
		}
		roomNumberHash := s.crypto.HashBasic(ctx, roomNumberBytes)
		repoFilter.RoomNumberHash = &roomNumberHash
	}

	count, err := s.roomRepo.Count(ctx, repoFilter)
	if err != nil {
		return 0, fmt.Errorf("count rooms: %w", err)
	}

	return count, nil
}


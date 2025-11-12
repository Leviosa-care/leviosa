package room

import (
	"context"
	"errors"
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
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid name value: %v", err))
		}
		nameHash := s.crypto.HashBasic(ctx, nameBytes)
		repoFilter.NameHash = &nameHash
	}

	if filter.RoomNumber != nil {
		roomNumberBytes, err := encx.SerializeValue(*filter.RoomNumber)
		if err != nil {
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid room number value: %v", err))
		}
		roomNumberHash := s.crypto.HashBasic(ctx, roomNumberBytes)
		repoFilter.RoomNumberHash = &roomNumberHash
	}

	roomsEncx, err := s.roomRepo.List(ctx, repoFilter)
	if err != nil {
		// Handle specific repository errors with context
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			return nil, fmt.Errorf("failed to list rooms due to database connection failure: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return nil, fmt.Errorf("failed to list rooms due to too many database connections: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, fmt.Errorf("failed to list rooms due to exhausted database resources: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("list rooms query was cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return nil, fmt.Errorf("failed to list rooms due to transaction failure: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			return nil, fmt.Errorf("failed to list rooms due to database deadlock: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, fmt.Errorf("failed to list rooms due to invalid filter parameters: %w", err)
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("list rooms operation was cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("list rooms operation timed out: %w", err)
		default:
			return nil, fmt.Errorf("failed to list rooms: %w", err)
		}
	}

	var rooms []*domain.RoomResponse
	for _, roomEncx := range roomsEncx {
		room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("room", err)
		}
		rooms = append(rooms, room.ToResponse())
	}

	return rooms, nil
}

package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *RoomService) CreateRoom(ctx context.Context, buildingID uuid.UUID, name string, capacity int) (*domain.Room, error) {
	// Verify building exists and is active
	building, err := s.buildingRepo.GetByID(ctx, buildingID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("verify building exists: %w", err)
	}

	if !building.IsActive {
		return nil, fmt.Errorf("cannot create room in inactive building")
	}

	// Create domain entity with validation
	room, err := domain.NewRoom(buildingID, name, capacity)
	if err != nil {
		return nil, fmt.Errorf("create room entity: %w", err)
	}

	// Persist to repository
	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	return room, nil
}

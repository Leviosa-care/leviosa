package room

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *RoomService) CreateRoom(ctx context.Context, request *domain.CreateRoomRequest) (*domain.RoomResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Verify building exists and is active
	buildingEncx, err := s.buildingRepo.GetByID(ctx, request.BuildingID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewInvalidInputErr(err)
		}
		return nil, fmt.Errorf("get building by ID: %w", err)
	}

	if !buildingEncx.IsActive {
		return nil, errs.NewInvalidValueErr("cannot create room in inactive building")
	}

	// Create domain entity
	now := time.Now()
	room := &domain.Room{
		ID:          uuid.New(),
		BuildingID:  request.BuildingID,
		Name:        request.Name,
		Description: request.Description,
		RoomNumber:  request.RoomNumber,
		Capacity:    request.Capacity,
		Equipment:   request.Equipment,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	roomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("room", err)
	}

	// Persist to repository
	if err := s.roomRepo.Create(ctx, roomEncx); err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	return room.ToResponse(), nil
}

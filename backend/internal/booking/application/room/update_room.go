package room

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// func (s *RoomService) UpdateRoom(ctx context.Context, id uuid.UUID, name, description, roomNumber string, capacity int) (*domain.RoomResponse, error) {
func (s *RoomService) UpdateRoom(ctx context.Context, request *domain.UpdateRoomRequest) (*domain.RoomResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing room
	roomEncx, err := s.roomRepo.GetByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("get room by ID: %w", err)
	}

	// Decrypt room
	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("room", err)
	}

	if request.Name != nil {
		room.Name = *request.Name
	}

	if request.Description != nil {
		room.Description = *request.Description
	}
	if request.RoomNumber != nil {
		room.RoomNumber = *request.RoomNumber
	}
	if request.Capacity != nil {
		room.Capacity = *request.Capacity
	}
	if request.Equipment != nil {
		room.Equipment = *request.Equipment
	}

	if request.IsActive != nil {
		room.IsActive = *request.IsActive
	}

	// Re-encrypt updated room
	updatedRoomEncx, err := domain.ProcessRoomEncx(ctx, s.crypto, room)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("room", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, updatedRoomEncx); err != nil {
		return nil, fmt.Errorf("update room: %w", err)
	}

	return room.ToResponse(), nil
}

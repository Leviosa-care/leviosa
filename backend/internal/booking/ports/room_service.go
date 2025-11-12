package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// RoomService defines the interface for room business logic
type RoomService interface {
	// CreateRoom creates a new room with validation
	CreateRoom(ctx context.Context, request *domain.CreateRoomRequest) (*domain.RoomResponse, error)

	// GetRoom retrieves a room by ID
	GetRoom(ctx context.Context, id uuid.UUID) (*domain.RoomResponse, error)

	//
	// // UpdateRoom updates room details with validation
	// UpdateRoom(ctx context.Context, id uuid.UUID, name, description, roomNumber string, capacity int) (*domain.RoomResponse, error)
	//
	// // SetRoomEquipment updates the room's equipment list
	// SetRoomEquipment(ctx context.Context, id uuid.UUID, equipment []string) (*domain.RoomResponse, error)
	//
	// // DeactivateRoom marks a room as inactive
	// DeactivateRoom(ctx context.Context, id uuid.UUID) error
	//
	// // ActivateRoom marks a room as active
	// ActivateRoom(ctx context.Context, id uuid.UUID) error
	//
	// // ListRooms retrieves rooms with filtering
	// ListRooms(ctx context.Context, filter RoomFilter) ([]*domain.Room, error)
	//
	// // GetRoomsByBuilding retrieves all rooms in a building
	// GetRoomsByBuilding(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.RoomResponse, error)
}

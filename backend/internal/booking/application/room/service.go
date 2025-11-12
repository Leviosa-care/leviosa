package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

type RoomService struct {
	roomRepo     ports.RoomRepository
	buildingRepo ports.BuildingRepository
	crypto       encx.CryptoService
}

// New creates a new instance of the room service
func New(roomRepo ports.RoomRepository, buildingRepo ports.BuildingRepository, crypto encx.CryptoService) ports.RoomService {
	return &RoomService{
		roomRepo:     roomRepo,
		buildingRepo: buildingRepo,
		crypto:       crypto,
	}
}

func (s *RoomService) GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get room: %w", err)
	}

	return room, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, id uuid.UUID, name, description, roomNumber string, capacity int) (*domain.Room, error) {
	// Get existing room
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get room for update: %w", err)
	}

	// Update details with validation
	if err := room.UpdateDetails(name, description, roomNumber, capacity); err != nil {
		return nil, fmt.Errorf("update room details: %w", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("update room: %w", err)
	}

	return room, nil
}

func (s *RoomService) SetRoomEquipment(ctx context.Context, id uuid.UUID, equipment []string) (*domain.Room, error) {
	// Get existing room
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get room for equipment update: %w", err)
	}

	// Set equipment
	room.SetEquipment(equipment)

	// Persist changes
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("update room equipment: %w", err)
	}

	return room, nil
}

func (s *RoomService) SetRoomRate(ctx context.Context, id uuid.UUID, rateCents int) (*domain.Room, error) {
	// Get existing room
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get room for rate update: %w", err)
	}

	// Set rate with validation
	if err := room.SetHourlyRate(rateCents); err != nil {
		return nil, fmt.Errorf("set room rate: %w", err)
	}

	// Persist changes
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("update room rate: %w", err)
	}

	return room, nil
}

func (s *RoomService) ClearRoomRate(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	// Get existing room
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get room for rate clearing: %w", err)
	}

	// Clear rate
	room.ClearHourlyRate()

	// Persist changes
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("clear room rate: %w", err)
	}

	return room, nil
}

func (s *RoomService) DeactivateRoom(ctx context.Context, id uuid.UUID) error {
	// Get existing room
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get room for deactivation: %w", err)
	}

	// Deactivate
	room.Deactivate()

	// Persist changes
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return fmt.Errorf("deactivate room: %w", err)
	}

	return nil
}

func (s *RoomService) ActivateRoom(ctx context.Context, id uuid.UUID) error {
	// Get existing room
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get room for activation: %w", err)
	}

	// Verify building is still active
	building, err := s.buildingRepo.GetByID(ctx, room.BuildingID)
	if err != nil {
		return fmt.Errorf("verify building for room activation: %w", err)
	}

	if !building.IsActive {
		return fmt.Errorf("cannot activate room in inactive building")
	}

	// Activate
	room.Activate()

	// Persist changes
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return fmt.Errorf("activate room: %w", err)
	}

	return nil
}

func (s *RoomService) ListRooms(ctx context.Context, filter ports.RoomFilter) ([]*domain.Room, error) {
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

	rooms, err := s.roomRepo.List(ctx, repoFilter)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid filter parameters: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during room list retrieval: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room list: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during room list retrieval: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during room list retrieval: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room list retrieval: %w", err))
		}
	}

	return rooms, nil
}

func (s *RoomService) GetRoomsByBuilding(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.Room, error) {
	// Verify building exists
	_, err := s.buildingRepo.GetByID(ctx, buildingID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("verify building exists: %w", err)
	}

	rooms, err := s.roomRepo.GetByBuildingID(ctx, buildingID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("get rooms by building: %w", err)
	}

	return rooms, nil
}

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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.ErrRepositoryNotFound
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error while verifying building: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("query failed while verifying building: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error while verifying building: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error while verifying building: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled error while verifying building: %w", err))
		}
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
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("room data: %v", err))
		case errors.Is(err, errs.ErrUniqueViolation):
			return nil, errs.NewAlreadyExistsError(err, "room with this name or room number")
		case errors.Is(err, errs.ErrNotNullViolation):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("missing required data for room: %v", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid foreign key for room: %v", err))
		case errors.Is(err, errs.ErrCheckViolation):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("room data failed check constraint: %v", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error for room: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during room creation: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room creation: %w", err))
		}
	}

	return room.ToResponse(), nil
}

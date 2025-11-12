package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *RoomService) GetRoom(ctx context.Context, id uuid.UUID) (*domain.RoomResponse, error) {
	roomEncx, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.ErrRepositoryNotFound
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid room ID: %v", err))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during room retrieval: %w", err))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for room: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database error during room retrieval: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during room retrieval: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during room retrieval: %w", err))
		}
	}

	room, err := domain.DecryptRoomEncx(ctx, s.crypto, roomEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("room", err)
	}

	return room.ToResponse(), nil
}

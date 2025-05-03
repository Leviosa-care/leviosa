package eventService

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// CreateEvent creates a new event in the system.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - event: A pointer to a models.Event instance representing the event to be created.
//
// Returns:
//   - string: The ID of the created event.
//   - error: An error if the event creation fails, the query fails, or an unexpected error occurs.
//     Returns nil if the event is created successfully.
func (s *service) CreateEvent(ctx context.Context, event *models.Event) (string, error) {
	day, month, year, err := ParseBeginAt(event)
	if err != nil {
		return "", domain.NewInvalidValueErr("invalid BeginAt")
	}

	if err = s.repo.IsDateAvailable(ctx, day, month, year); err != nil {
		switch {
		case errors.Is(err, rp.ErrValidation):
			return "", domain.NewNotCreatedErr(err)
		case errors.Is(err, rp.ErrContext):
			return "", err
		case errors.Is(err, rp.ErrDatabase):
			return "", domain.NewQueryFailedErr(err)
		}
	}
	event.Day = day
	event.Month = month
	event.Year = year

	event.ID = uuid.NewString()

	if err := s.crypto.ProcessStruct(ctx, event); err != nil {
		return "", domain.NewNotEncryptedErr("event", err)
	}

	eventID, err := s.repo.AddEvent(ctx, event)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotCreated):
			return "", domain.NewNotCreatedErr(err)
		case errors.Is(err, rp.ErrDatabase):
			return "", domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return "", err
		}
	}
	return eventID, nil
}

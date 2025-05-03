package eventService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) ModifyEvent(ctx context.Context, event *models.Event) error {
	if event == nil {
		return domain.NewInvalidValueErr("event can not be nil")
	}
	whereMap := map[string]any{"id": event.ID}
	if err := s.repo.ModifyEvent(ctx,
		event,
		whereMap,
	); err != nil {
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrInternal):
			return err
		case errors.Is(err, rp.ErrNotUpdated):
			return domain.NewNotUpdatedErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}

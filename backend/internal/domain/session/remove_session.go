package sessionService

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) RemoveSession(ctx context.Context, sessionID string) error {
	if err := uuid.Validate(sessionID); err != nil {
		return domain.NewInvalidValueErr(fmt.Sprintf("invalid sessionID: %s", err))
	}
	err := s.Repo.RemoveSession(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrNotFound):
			return domain.NewNotFoundErr(err)
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}

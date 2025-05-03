package throttlerService

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) Reset(ctx context.Context, email string) error {
	throttlerEncoded, err := s.repo.IsLocked(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		}
	}
	var res Info
	if err = json.Unmarshal(throttlerEncoded, &res); err != nil {
		return domain.NewJSONUnmarshalErr(err)
	}

	if time.Now().Before(res.LockedUntil) {
		return domain.NewLockedAccountErr(err, "throttler")
	}

	if err := s.repo.Reset(ctx, email); err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return domain.NewNotFoundErr(err)
		default:
			return domain.NewQueryFailedErr(err)
		}
	}
	return nil
}

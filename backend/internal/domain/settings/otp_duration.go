package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetOTPDuration(ctx context.Context) (int, error) {
	setting, err := s.repo.GetInt(ctx, OTPDurationKey)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return 0, domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return 0, err
		case errors.Is(err, rp.ErrDatabase):
			return 0, domain.NewQueryFailedErr(err)
		}
	}
	return setting.Value, nil
}

func (s *service) SetOTPDuration(ctx context.Context, duration int) error {
	setting := NewSetting(OTPDurationKey, duration)
	return s.repo.SetInt(ctx, setting)
}

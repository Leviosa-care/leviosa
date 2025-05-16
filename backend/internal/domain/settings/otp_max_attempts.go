package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetOTPMaxAttempts(ctx context.Context) (int, error) {
	setting, err := s.repo.GetInt(ctx, OTPMaxAttemptsKey)
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

func (s *service) SetOTPMaxAttempts(ctx context.Context, value int) error {
	setting := NewSetting(OTPMaxAttemptsKey, value)
	return s.repo.SetInt(ctx, setting)
}

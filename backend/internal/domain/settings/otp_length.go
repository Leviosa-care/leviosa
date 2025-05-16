package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetOTPLength(ctx context.Context) (int, error) {
	setting, err := s.repo.GetInt(ctx, OTPLengthKey)
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

func (s *service) SetOTPLength(ctx context.Context, length int) error {
	setting := NewSetting(OTPLengthKey, length)
	return s.repo.SetInt(ctx, setting)
}

package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetOTPDuration(ctx context.Context) (*domain.GetOTPDurationResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.OTPDuration)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "OTP duration")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetOTPDurationResponse{Duration: setting.Value}, nil
}

func (s *SettingsService) SetOTPDuration(ctx context.Context, request *domain.SetOTPDurationRequest) (*domain.SetOTPDurationResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.OTPDuration, request.Duration)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.OTPDuration, request.Duration); err != nil {
		return nil, err
	}

	return &domain.SetOTPDurationResponse{Success: true}, nil
}

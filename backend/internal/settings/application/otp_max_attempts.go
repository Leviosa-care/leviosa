package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetOTPMaxAttempts(ctx context.Context) (*domain.GetOTPMaxAttemptsResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.OTPMaxAttempts)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "OTP max attempts")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetOTPMaxAttemptsResponse{MaxAttempts: setting.Value}, nil
}

func (s *SettingsService) SetOTPMaxAttempts(ctx context.Context, request *domain.SetOTPMaxAttemptsRequest) (*domain.SetOTPMaxAttemptsResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.OTPMaxAttempts, request.MaxAttempts)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, err
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.OTPMaxAttempts, request.MaxAttempts); err != nil {
	// 	return nil, err
	// }

	return &domain.SetOTPMaxAttemptsResponse{Success: true}, nil
}

package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetOTPLength(ctx context.Context) (*domain.GetOTPLengthResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.OTPLength)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "OTP length")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetOTPLengthResponse{Length: setting.Value}, nil
}

func (s *SettingsService) SetOTPLength(ctx context.Context, request *domain.SetOTPLengthRequest) (*domain.SetOTPLengthResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.OTPLength, request.Length)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, err
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.OTPLength, request.Length); err != nil {
	// 	return nil, err
	// }

	return &domain.SetOTPLengthResponse{Success: true}, nil
}

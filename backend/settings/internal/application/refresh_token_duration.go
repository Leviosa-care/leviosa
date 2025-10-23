package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetRefreshTokenDuration(ctx context.Context) (*domain.GetRefreshTokenDurationResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.RefreshTokenDuration)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "refresh token duration")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetRefreshTokenDurationResponse{Duration: setting.Value}, nil
}

func (s *SettingsService) SetRefreshTokenDuration(ctx context.Context, request *domain.SetRefreshTokenDurationRequest) (*domain.SetRefreshTokenDurationResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.RefreshTokenDuration, request.Duration)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.RefreshTokenDuration, request.Duration); err != nil {
		return nil, err
	}

	return &domain.SetRefreshTokenDurationResponse{Success: true}, nil
}
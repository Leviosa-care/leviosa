package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetAccessTokenDuration(ctx context.Context) (*domain.GetAccessTokenDurationResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.AccessTokenDuration)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "access token duration")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetAccessTokenDurationResponse{Duration: setting.Value}, nil
}

func (s *SettingsService) SetAccessTokenDuration(ctx context.Context, request *domain.SetAccessTokenDurationRequest) (*domain.SetAccessTokenDurationResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.AccessTokenDuration, request.Duration)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.AccessTokenDuration, request.Duration); err != nil {
		return nil, err
	}

	return &domain.SetAccessTokenDurationResponse{Success: true}, nil
}


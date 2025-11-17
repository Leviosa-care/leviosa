package settings

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetRefreshTokenDuration(ctx context.Context) (*domain.GetRefreshTokenDurationResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("get refresh token duration: %w", err)
	}
	return &domain.GetRefreshTokenDurationResponse{Duration: setting.Value}, nil
}

func (s *SettingsService) SetRefreshTokenDuration(ctx context.Context, request *domain.SetRefreshTokenDurationRequest) (*domain.SetRefreshTokenDurationResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.RefreshTokenDuration, request.Duration)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, fmt.Errorf("set refresh token duration: %w", err)
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.RefreshTokenDuration, request.Duration); err != nil {
	// 	return nil, err
	// }

	return &domain.SetRefreshTokenDurationResponse{Success: true}, nil
}


package settings

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetAccessTokenDuration(ctx context.Context) (*domain.GetAccessTokenDurationResponse, error) {
	setting, err := s.repo.GetInt(ctx, settings.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("get access token duration: %w", err)
	}
	return &domain.GetAccessTokenDurationResponse{Duration: setting.Value}, nil
}

func (s *SettingsService) SetAccessTokenDuration(ctx context.Context, request *domain.SetAccessTokenDurationRequest) (*domain.SetAccessTokenDurationResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.AccessTokenDuration, request.Duration)
	if err := s.repo.SetInt(ctx, setting); err != nil {
		return nil, fmt.Errorf("set access token duration: %w", err)
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.AccessTokenDuration, request.Duration); err != nil {
	// 	return nil, err
	// }

	return &domain.SetAccessTokenDurationResponse{Success: true}, nil
}

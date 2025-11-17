package settings

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetCompanyInstagram(ctx context.Context) (*domain.GetCompanyInstagramResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyInstagram)
	if err != nil {
		return nil, fmt.Errorf("get company instagram: %w", err)
	}
	return &domain.GetCompanyInstagramResponse{Instagram: setting.Value}, nil
}

func (s *SettingsService) SetCompanyInstagram(ctx context.Context, request *domain.SetCompanyInstagramRequest) (*domain.SetCompanyInstagramResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.CompanyInstagram, request.Instagram)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return nil, fmt.Errorf("set company instagram: %w", err)
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.CompanyInstagram, request.Instagram); err != nil {
	// 	return nil, err
	// }

	return &domain.SetCompanyInstagramResponse{Success: true}, nil
}

package settings

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func (s *SettingsService) GetCompanyName(ctx context.Context) (*domain.GetCompanyNameResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyName)
	if err != nil {
		return nil, fmt.Errorf("get company name: %w", err)
	}
	return &domain.GetCompanyNameResponse{Name: setting.Value}, nil
}

func (s *SettingsService) SetCompanyName(ctx context.Context, request *domain.SetCompanyNameRequest) (*domain.SetCompanyNameResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.CompanyName, request.Name)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return nil, err
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.CompanyName, request.Name); err != nil {
	// 	return nil, err
	// }

	return &domain.SetCompanyNameResponse{Success: true}, nil
}

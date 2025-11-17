package settings

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func (s *SettingsService) GetCompanyEmail(ctx context.Context) (*domain.GetCompanyEmailResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyEmail)
	if err != nil {
		return nil, fmt.Errorf("get company email: %w", err)
	}

	return &domain.GetCompanyEmailResponse{Email: setting.Value}, nil
}

func (s *SettingsService) SetCompanyEmail(ctx context.Context, request *domain.SetCompanyEmailRequest) (*domain.SetCompanyEmailResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.CompanyEmail, request.Email)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return nil, err
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.CompanyEmail, request.Email); err != nil {
	// 	return nil, err
	// }

	return &domain.SetCompanyEmailResponse{Success: true}, nil
}

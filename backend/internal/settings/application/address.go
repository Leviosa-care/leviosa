package settings

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetCompanyLegalAddress(ctx context.Context) (*domain.GetCompanyLegalAddressResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyLegalAddress)
	if err != nil {
		return nil, fmt.Errorf("get company legal address: %w", err)
	}
	return &domain.GetCompanyLegalAddressResponse{Address: setting.Value}, nil
}

func (s *SettingsService) SetCompanyLegalAddress(ctx context.Context, request *domain.SetCompanyLegalAddressRequest) (*domain.SetCompanyLegalAddressResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.CompanyLegalAddress, request.Address)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return nil, fmt.Errorf("set company legal address: %w", err)
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.CompanyLegalAddress, request.Address); err != nil {
	// 	return nil, err
	// }

	return &domain.SetCompanyLegalAddressResponse{Success: true}, nil
}

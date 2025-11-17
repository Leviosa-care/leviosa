package settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"
)

func (s *SettingsService) GetCompanyTelephone(ctx context.Context) (*domain.GetCompanyTelephoneResponse, error) {
	settingEncx, err := s.repo.GetEncryptedSetting(ctx, settings.CompanyPhone)
	if err != nil {
		return nil, fmt.Errorf("get company telephone: %w", err)
	}

	// Use generated decrypt function
	setting, err := domain.DecryptSettingEncryptedEncx(ctx, s.crypto, settingEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("settings phone number", err)
	}

	return &domain.GetCompanyTelephoneResponse{Telephone: setting.Value}, nil
}

func (s *SettingsService) SetCompanyTelephone(ctx context.Context, request *domain.SetCompanyTelephoneRequest) (*domain.SetCompanyTelephoneResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Trim whitespace before processing
	trimmedPhone := strings.TrimSpace(request.Telephone)

	if err := validation.ValidatePhone(trimmedPhone); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Create clean source struct
	setting := &domain.SettingEncrypted{
		Key:   settings.CompanyPhone,
		Value: trimmedPhone,
	}

	// Use generated process function - handles DEK generation automatically
	settingEncx, err := domain.ProcessSettingEncryptedEncx(ctx, s.crypto, setting)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("settings phone number", err)
	}

	if err := s.repo.SetEncryptedSetting(ctx, settingEncx); err != nil {
		return nil, fmt.Errorf("set company telephone: %w", err)
	}

	// COMMENTED OUT: Event publishing disabled - other modules will access settings via interface
	// See CLAUDE.local.md for details on the new architecture pattern
	// if err := s.publisher.PublishSettingUpdate(ctx, settings.CompanyPhone, trimmedPhone); err != nil {
	// 	return nil, err
	// }

	return &domain.SetCompanyTelephoneResponse{Success: true}, nil
}

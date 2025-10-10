package settings

import (
	"context"
	"errors"
	"strings"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/validation"
)

func (s *SettingsService) GetCompanyTelephone(ctx context.Context) (*domain.GetCompanyTelephoneResponse, error) {
	settingEncx, err := s.repo.GetEncryptedSetting(ctx, settings.CompanyPhone)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "company telephone")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
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
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.CompanyPhone, trimmedPhone); err != nil {
		return nil, err
	}

	return &domain.SetCompanyTelephoneResponse{Success: true}, nil
}

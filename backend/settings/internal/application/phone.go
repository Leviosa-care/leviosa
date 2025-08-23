package settings

import (
	"context"
	"errors"
	"strings"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
)

func (s *SettingsService) GetCompanyTelephone(ctx context.Context) (*domain.GetCompanyTelephoneResponse, error) {
	setting, err := s.repo.GetEncryptedSetting(ctx, settings.CompanyPhone)
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
	if err := s.crypto.DecryptStruct(ctx, setting); err != nil {
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

	if err := domain.ValidateTelephone(trimmedPhone); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting, err := NewSettingEncrypted(s, settings.CompanyPhone, trimmedPhone)
	if err != nil {
		return nil, errs.NewNotCreatedErr(err, "company telephone")
	}
	if err := s.crypto.ProcessStruct(ctx, setting); err != nil {
		return nil, errs.NewNotEncryptedErr("settings phone number", err)
	}

	if err := s.repo.SetEncryptedSetting(ctx, setting); err != nil {
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.CompanyPhone, trimmedPhone); err != nil {
		return nil, err
	}

	return &domain.SetCompanyTelephoneResponse{Success: true}, nil
}

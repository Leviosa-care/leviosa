package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetCompanyLegalAddress(ctx context.Context) (*domain.GetCompanyLegalAddressResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyLegalAddress)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "company legal address")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetCompanyLegalAddressResponse{Address: setting.Value}, nil
}

func (s *SettingsService) SetCompanyLegalAddress(ctx context.Context, request *domain.SetCompanyLegalAddressRequest) (*domain.SetCompanyLegalAddressResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.CompanyLegalAddress, request.Address)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.CompanyLegalAddress, request.Address); err != nil {
		return nil, err
	}

	return &domain.SetCompanyLegalAddressResponse{Success: true}, nil
}

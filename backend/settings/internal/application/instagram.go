package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetCompanyInstagram(ctx context.Context) (*domain.GetCompanyInstagramResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyInstagram)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "company instagram")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
	}
	return &domain.GetCompanyInstagramResponse{Instagram: setting.Value}, nil
}

func (s *SettingsService) SetCompanyInstagram(ctx context.Context, request *domain.SetCompanyInstagramRequest) (*domain.SetCompanyInstagramResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	setting := NewSetting(settings.CompanyInstagram, request.Instagram)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return nil, err
	}

	if err := s.PublishSettingUpdate(ctx, settings.CompanyInstagram, request.Instagram); err != nil {
		return nil, err
	}

	return &domain.SetCompanyInstagramResponse{Success: true}, nil
}

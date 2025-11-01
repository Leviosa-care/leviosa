package settings

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *SettingsService) GetCompanyEmail(ctx context.Context) (*domain.GetCompanyEmailResponse, error) {
	setting, err := s.repo.GetString(ctx, settings.CompanyEmail)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "company email")
		case errors.Is(err, errs.ErrContext):
			return nil, err
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewQueryFailedErr(err)
		}
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

package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetCompanyEmail(ctx context.Context) (string, error) {
	setting, err := s.repo.GetString(ctx, CompanyEmailKey)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return "", domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return "", err
		case errors.Is(err, rp.ErrDatabase):
			return "", domain.NewQueryFailedErr(err)
		}
	}

	return setting.Value, nil
}

func (s *service) SetCompanyEmail(ctx context.Context, email string) error {
	setting := NewSetting(CompanyEmailKey, email)
	if err := s.repo.SetString(ctx, setting); err != nil {
		return err
	}
	return s.PublishSettingUpdate(ctx, CompanyEmailKey, email)
}

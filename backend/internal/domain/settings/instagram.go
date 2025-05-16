package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetCompanyInstagram(ctx context.Context) (string, error) {
	setting, err := s.repo.GetString(ctx, CompanyInstagramKey)
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

func (s *service) SetCompanyInstagram(ctx context.Context, link string) error {
	setting := NewSetting(CompanyInstagramKey, link)
	return s.repo.SetString(ctx, setting)
}

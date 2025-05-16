package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetCompanyLegalAddress(ctx context.Context) (string, error) {
	setting, err := s.repo.GetString(ctx, CompanyLegalAddressKey)
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

func (s *service) SetCompanyLegalAddress(ctx context.Context, address string) error {
	setting := NewSetting(CompanyLegalAddressKey, address)
	return s.repo.SetString(ctx, setting)
}

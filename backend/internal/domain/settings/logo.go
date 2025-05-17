package settings

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) SetCompanyLogo(ctx context.Context, logo []byte) error {
	if err := s.media.SetLogo(ctx, logo); err != nil {
		return err
	}
	return s.PublishSettingUpdate(ctx, CompanyEmailKey, logo)
}

func (s *service) GetCompanyLogo(ctx context.Context) ([]byte, error) {
	logo, err := s.media.GetLogo(ctx)
	if errors.Is(err, rp.ErrNotFound) {
		return nil, domain.NewNotFoundErr(err)
	}
	return logo, nil
}

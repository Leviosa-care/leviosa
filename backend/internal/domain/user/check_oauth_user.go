package userService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) CheckOAuthUser(ctx context.Context, email string, provider models.ProviderType) error {
	if err := models.ValidateEmail(email); err != nil {
		return domain.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, []byte(email))
	if err := s.repo.HasOAuthUser(ctx, emailHash, provider); err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}

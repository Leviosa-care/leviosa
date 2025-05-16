package otpService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) CancelOTP(ctx context.Context, email string) error {
	emailHash := s.crypto.HashBasic(ctx, []byte(email))
	err := s.repo.InvalidateOTP(ctx, emailHash)
	switch {
	case errors.Is(err, rp.ErrNotFound):
		switch {
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrNotDeleted):
			return domain.NewNotDeletedErr(err)
		}
	}
	return nil
}

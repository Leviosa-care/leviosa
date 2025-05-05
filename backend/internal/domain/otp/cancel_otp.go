package otpService

import (
	"context"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) CancelOTP(ctx context.Context, email string) error {
	emailHash := s.crypto.HashBasic(ctx, []byte(email))
	err := s.Repo.InvalidateOTP(ctx, emailHash)
	switch {
	case errors.Is(err, rp.ErrNotFound):
		// TODO: change the error returned here brother
		return nil
	}
	return nil
}

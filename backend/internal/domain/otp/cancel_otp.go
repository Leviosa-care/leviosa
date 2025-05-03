package otpService

import (
	"context"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) CancelOTP(ctx context.Context, email string) error {
	err := s.Repo.InvalidateOTP(ctx, email)
	switch {
	case errors.Is(err, rp.ErrNotFound):
		// TODO: change the error returned here brother
		return nil
	}
	return nil
}

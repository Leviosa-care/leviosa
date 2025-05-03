package otpService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) ValidateOTP(ctx context.Context, emailHash string, value string) error {
	// OTPencoded, err := s.Repo.GetOTP(ctx, otp.Email)
	err := s.Repo.ValidateOTP(ctx, emailHash, value)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrValidation):
			return domain.NewInvalidValueErr("invalid OTP")
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}

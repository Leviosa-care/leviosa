package otp

import (
	"context"
	"errors"

	"github.com/Leviosa-care/core/errs"
)

func (s *OTPService) CancelOTP(ctx context.Context, email string) error {
	if email == "" {
		return errs.NewInvalidValueErr("email is required")
	}

	// Hash email for lookup
	emailHash := s.crypto.HashBasic(ctx, []byte(email))

	// Attempt to invalidate the OTP
	err := s.repo.InvalidateOTP(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "OTP")
		case errors.Is(err, errs.ErrContext):
			return err // Pass through context errors
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewUnexpectedError(err) // Connection/network issues
		default:
			return errs.NewNotDeletedErr(err, "OTP")
		}
	}

	return nil
}

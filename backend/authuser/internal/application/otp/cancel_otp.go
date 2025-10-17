package otp

import (
	"context"
	"errors"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/validation"
	"github.com/hengadev/encx"
)

func (s *OTPService) CancelOTP(ctx context.Context, email string) error {
	if err := validation.ValidateEmail(email); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Hash email for lookup
	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	// Attempt to invalidate the OTP
	if err = s.repo.InvalidateOTP(ctx, emailHash); err != nil {
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

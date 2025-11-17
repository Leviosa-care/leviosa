package otp

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"
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
		return fmt.Errorf("invalidate OTP: %w", err)
	}

	return nil
}

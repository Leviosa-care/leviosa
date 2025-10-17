package otp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/contracts/otp"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/validation"
	"github.com/hengadev/encx"
)

func (s *OTPService) ResendOTP(ctx context.Context, email string) error {
	if err := validation.ValidateEmail(email); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	// Check if existing OTP exists
	marshaledOTP, err := s.repo.GetOTP(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "OTP")
		default:
			return errs.NewUnexpectedError(err)
		}
	}

	// Deserialize and decrypt existing OTP using the new generated function
	var existingOTPEncx domain.OTPEncx
	if err := json.Unmarshal(marshaledOTP, &existingOTPEncx); err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	existingOTP, err := domain.DecryptOTPEncx(ctx, s.crypto, &existingOTPEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("OTP", err)
	}

	// Check if OTP has expired
	if existingOTP.IsExpired() {
		// Clean up expired OTP
		if err := s.repo.InvalidateOTP(ctx, emailHash); err != nil {
			switch {
			case errors.Is(err, errs.ErrRepositoryNotFound):
				// Already cleaned up - that's fine
			case errors.Is(err, errs.ErrContext):
				return err
			case errors.Is(err, errs.ErrDBQuery):
				return errs.NewUnexpectedError(err)
			default:
				// Log error but continue
				// TODO: Add proper logging
			}
		}
		return errs.NewExpiredTokenErr("OTP", errors.New("OTP has expired"))
	}

	// Check if max attempts exceeded
	if existingOTP.Attempts >= s.GetOTPMaxAttempts() {
		return errs.NewRateLimitErr(errors.New("maximum attempts exceeded"), "OTP")
	}

	if err := s.PublishOTPUpdate(
		ctx,
		otp.Email,
		&domain.OTPSentEvent{
			Code:      existingOTP.Code,
			Email:     existingOTP.Email,
			ExpiresAt: existingOTP.ExpiresAt,
		},
	); err != nil {
		return errs.NewExternalServiceErr(err, "publish OTP update")
	}

	return nil
}

package otp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (s *OTPService) CreateOTP(ctx context.Context, email string) error {
	if email == "" {
		return errs.NewInvalidValueErr("email is required")
	}

	emailHash := s.crypto.HashBasic(ctx, []byte(email))

	// Check if OTP already exists and is still valid
	marshaledOTP, err := s.repo.GetOTP(ctx, emailHash)
	if err != nil {
		if !errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewUnexpectedError(err)
		}
		// No existing OTP found - continue to create new one
	} else {
		// OTP exists, check if still valid
		var existingOTP domain.OTP
		if err := json.Unmarshal(marshaledOTP, &existingOTP); err != nil {
			return errs.NewJSONUnmarshalErr(err)
		}

		if err := s.crypto.DecryptStruct(ctx, &existingOTP); err != nil {
			return errs.NewNotDecryptedErr("OTP", err)
		}

		if !existingOTP.IsExpired() && existingOTP.Attempts < s.GetOTPMaxAttempts() {
			return errs.NewRateLimitErr(errors.New("OTP already active"), "OTP")
		}
	}

	// Generate new OTP
	otp, err := s.generateOTP(email)
	if err != nil {
		return errs.NewNotCreatedErr(err, "OTP")
	}

	// Encrypt OTP data
	if err := s.crypto.ProcessStruct(ctx, otp); err != nil {
		return errs.NewNotEncryptedErr("OTP", err)
	}
	// Serialize and save to repository
	otpData, err := json.Marshal(otp)
	if err != nil {
		return errs.NewJSONMarshalErr(err)
	}

	ttl := time.Duration(s.GetOTPDuration()) * time.Minute
	if err := s.repo.SaveOTP(ctx, otp.EmailHash, otpData, ttl); err != nil {
		switch {
		case errors.Is(err, errs.ErrContext):
			return err // Pass through context errors
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewUnexpectedError(err) // Connection/network issues
		default:
			return errs.NewNotCreatedErr(err, "OTP")
		}
	}

	if err := s.PublishOTPUpdate(
		ctx,
		otp.Email,
		&domain.OTPSentEvent{
			Code:      otp.Code,
			Email:     otp.Email,
			ExpiresAt: otp.ExpiresAt,
		},
	); err != nil {
		return errs.NewExternalServiceErr(err, "publish OTP update")
	}

	return nil
}

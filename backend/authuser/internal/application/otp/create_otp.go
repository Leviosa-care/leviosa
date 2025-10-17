package otp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/hengadev/encx"
)

func (s *OTPService) CreateOTP(ctx context.Context, email string) error {
	if email == "" {
		return errs.NewInvalidValueErr("email is required")
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	// Check if OTP already exists and is still valid
	marshaledOTP, err := s.repo.GetOTP(ctx, emailHash)
	if err != nil {
		if !errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewUnexpectedError(err)
		}
		// No existing OTP found - continue to create new one
	} else {
		// OTP exists, check if still valid
		var existingOTPEncx domain.OTPEncx
		if err := json.Unmarshal(marshaledOTP, &existingOTPEncx); err != nil {
			return errs.NewJSONUnmarshalErr(err)
		}

		existingOTP, err := domain.DecryptOTPEncx(ctx, s.crypto, &existingOTPEncx)
		if err != nil {
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
	otpEncx, err := domain.ProcessOTPEncx(ctx, s.crypto, otp)
	if err != nil {
		return errs.NewNotEncryptedErr("OTP", err)
	}
	// Serialize and save to repository
	otpData, err := json.Marshal(otpEncx)
	if err != nil {
		return errs.NewJSONMarshalErr(err)
	}

	ttl := time.Duration(s.GetOTPDuration()) * time.Minute
	if err := s.repo.SaveOTP(ctx, otpEncx.EmailHash, otpData, ttl); err != nil {
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

package otp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/hengadev/encx"
)

func (s *OTPService) RequestOTP(ctx context.Context, email string) error {
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// No existing OTP found - continue to create new one
		case errors.Is(err, errs.ErrConnectionFailure):
			// Redis connection issues
			return fmt.Errorf("get existing OTP: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			// Redis connection pool exhausted
			return fmt.Errorf("get existing OTP: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			// Redis operation timeout
			return fmt.Errorf("get existing OTP: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			// Redis memory issues
			return fmt.Errorf("get existing OTP: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			// General Redis error
			return fmt.Errorf("get existing OTP: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return fmt.Errorf("get existing OTP cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return fmt.Errorf("get existing OTP timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return fmt.Errorf("get existing OTP: %w", err)
		}
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

		if !existingOTP.IsExpired() && existingOTP.Attempts < defaultOTPMaxAttempts {
			return errs.NewRateLimitErr(errors.New("OTP already active"), "OTP")
		}
	}

	// Generate new OTP
	otp, err := s.generateOTP(email)
	if err != nil {
		return errs.NewNotCreatedErr(err, "OTP")
	}

	// Encrypt OTP data using the new generated function
	otpEncx, err := domain.ProcessOTPEncx(ctx, s.crypto, otp)
	if err != nil {
		return errs.NewNotEncryptedErr("OTP", err)
	}
	// Serialize and save to repository
	otpData, err := json.Marshal(otpEncx)
	if err != nil {
		return errs.NewJSONMarshalErr(err)
	}

	ttl := time.Duration(defaultOTPDuration) * time.Minute

	if err := s.repo.SaveOTP(ctx, otpEncx.EmailHash, otpData, ttl); err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure):
			// Redis connection issues
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			// Redis connection pool exhausted
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			// Redis operation timeout
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			// Redis transaction failed
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			// Redis memory issues
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			// Redis authentication issues
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			// General Redis error
			return fmt.Errorf("save OTP: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return fmt.Errorf("save OTP cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return fmt.Errorf("save OTP timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return fmt.Errorf("save OTP: %w", err)
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

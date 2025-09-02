package otp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (s *OTPService) RequestOTP(ctx context.Context, email string) error {
	if email == "" {
		return errs.NewInvalidValueErr("email is required")
	}

	emailHash := s.crypto.HashBasic(ctx, []byte(email))

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

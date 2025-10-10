package otp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
)

func (s *OTPService) ValidateOTP(ctx context.Context, request *domain.ValidateOTPRequest) error {
	if err := request.Valid(ctx, s.effectiveOTPLength()); err != nil {
		return errs.NewInvalidValueErr(fmt.Sprintf("OTP validation failed: %s", err.Error()))
	}

	// Create a temporary OTP to get the email hash (since there's no standalone hash function)
	tempOTP := &domain.OTP{
		Email: request.Email,
	}
	tempOTPEncx, err := domain.ProcessOTPEncx(ctx, s.crypto, tempOTP)
	if err != nil {
		return errs.NewNotEncryptedErr("temporary OTP for email hash", err)
	}
	emailHash := tempOTPEncx.EmailHash

	// Get OTP from repository
	otpData, err := s.repo.GetOTP(ctx, emailHash)
	if err != nil {
		return s.handleRepositoryError(err, "retrieve OTP")
	}

	// Deserialize OTP
	var otpEncx domain.OTPEncx
	if err := json.Unmarshal(otpData, &otpEncx); err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	// Decrypt OTP using the new generated function
	otp, err := domain.DecryptOTPEncx(ctx, s.crypto, &otpEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("OTP", err)
	}

	// Check if OTP is expired
	if otp.IsExpired() {
		// Clean up expired OTP (best effort - ignore errors)
		s.cleanupOTP(ctx, emailHash, "expired OTP cleanup")
		return errs.NewExpiredTokenErr("OTP", errors.New("OTP has expired"))
	}

	// Check if max attempts exceeded
	if otp.Attempts >= s.GetOTPMaxAttempts() {
		// Clean up OTP that exceeded attempts (best effort - ignore errors)
		s.cleanupOTP(ctx, emailHash, "max attempts exceeded cleanup")
		return errs.NewRateLimitErr(errors.New("maximum attempts exceeded"), "OTP")
	}

	// Verify the code
	if request.Code != otp.Code {
		// Increment attempts and save back
		otp.IncrementAttempts()

		// Calculate remaining TTL with buffer to handle clock skew
		remainingTTL := time.Until(otp.ExpiresAt)
		if remainingTTL <= 30*time.Second {
			// Too close to expiry, don't save and treat as expired (best effort - ignore errors)
			s.cleanupOTP(ctx, emailHash, "near expiry cleanup")
			return errs.NewExpiredTokenErr("OTP", errors.New("OTP has expired"))
		}

		// Re-encrypt and save updated OTP using the new generated function
		updatedOTPEncx, err := domain.ProcessOTPEncx(ctx, s.crypto, otp)
		if err != nil {
			return errs.NewNotEncryptedErr("OTP", err)
		}

		updatedData, err := json.Marshal(&updatedOTPEncx)
		if err != nil {
			return errs.NewJSONMarshalErr(err)
		}

		// Save updated OTP with incremented attempts and remaining TTL
		if err := s.repo.SaveOTP(ctx, emailHash, updatedData, remainingTTL); err != nil {
			return s.handleRepositoryError(err, "update OTP attempts")
		}

		return errs.NewValueMismatchErr(otp.Code, request.Code)
	}

	// OTP is valid - try to claim it exclusively by deleting it
	if err := s.cleanupOTP(ctx, emailHash, "successful validation cleanup"); err != nil {
		return err // Another request already consumed this OTP
	}

	return nil
}

// handleRepositoryError centralizes repository error handling
func (s *OTPService) handleRepositoryError(err error, operation string) error {
	switch {
	case errors.Is(err, errs.ErrRepositoryNotFound):
		return errs.NewNotFoundErr(err, "OTP")
	case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
		return errs.NewExternalServiceErr(err, "database unavailable")
	case errors.Is(err, errs.ErrInvalidInput):
		return errs.NewInvalidValueErr("invalid OTP data")
	default:
		return errs.NewInternalErr(fmt.Errorf("failed to %s: %w", operation, err))
	}
}

// cleanupOTP performs OTP cleanup and returns error if already consumed by concurrent request
func (s *OTPService) cleanupOTP(ctx context.Context, emailHash string, reason string) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		fmt.Printf("failed to retrieve logger in cleanupOTP: %v\n", err)
		// Continue without logging if logger retrieval fails
	}

	if err := s.repo.InvalidateOTP(ctx, emailHash); err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// OTP was already deleted by another concurrent request
			return errs.NewAlreadyConsumedErr("OTP")
		}

		// For other errors, log but don't fail the validation (best effort)
		if logger != nil {
			logger.WarnContext(ctx, "OTP cleanup failed",
				"reason", reason,
				"email_hash", emailHash,
				"error", err)
		} else {
			// Fallback logging if logger is not available
			fmt.Printf("OTP cleanup failed: reason=%s, email_hash=%s, error=%v\n",
				reason, emailHash, err)
		}

		// Don't fail validation for infrastructure errors
		return nil
	}

	return nil // Successfully cleaned up OTP
}

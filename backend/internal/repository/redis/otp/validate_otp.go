package otpRepository

import (
	"context"
	"errors"
	"time"

	"github.com/hengadev/leviosa/internal/domain/otp"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// ValidateOTP validates an OTP for a given email hash.
// It checks if the OTP exists, ensures it is not expired, verifies it hasn't exceeded the maximum attempts (deleting it if so),
// and removes the OTP if it is successfully validated.
func (o *Repository) ValidateOTP(ctx context.Context, emailHash, providedOTP string) error {
	key := formatOTPKey(emailHash)
	otpData, err := o.getExistingOTP(ctx, key)
	if err != nil {
		return err
	}

	if otpData == nil {
		return rp.NewNotFoundErr(errors.New("OTP not found or expired"), "otp")
	}

	if time.Now().After(otpData.ExpiresAt) {
		return rp.NewValidationErr(errors.New("expired OTP"), "OTP")
	}

	// Check attempts
	if otpData.Attempts >= otpService.MaxOTPAttempts {
		// Delete expired OTP
		o.client.Del(ctx, key)
		return rp.NewValidationErr(errors.New("max attempts exceeded"), "OTP")
	}

	// Increment attempts
	otpData.Attempts++
	if err := o.setOTP(ctx, key, otpData); err != nil {
		return err
	}

	// Validate OTP
	if providedOTP != otpData.Code {
		return rp.NewValidationErr(errors.New("provided OTP code does not match stored OTP code"), "OTP")
	}

	// Delete OTP after successful validation
	o.client.Del(ctx, key)
	return nil
}

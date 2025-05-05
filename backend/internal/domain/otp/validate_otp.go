package otpService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) ValidateOTP(ctx context.Context, email string, value string) error {
	// create email hash
	emailHash := s.crypto.HashBasic(ctx, []byte(email))
	// get otp
	otpEncoded, err := s.Repo.GetOTP(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrNotFound):
			// TODO: return some error because the OTP is not found for the email
			return domain.NewNotFoundErr(err)
		}
	}
	// decode the value (deserialization)
	var data OTPData
	if err := json.Unmarshal(otpEncoded, &data); err != nil {
		return rp.NewDatabaseErr(fmt.Errorf("failed to parse existing OTP data: %w", err))
	}
	// validation logic
	if time.Now().After(data.ExpiresAt) {
		return domain.NewInvalidValueErr("expired OTP")
	}
	if data.Attempts >= MaxOTPAttempts {
		// delete expired OTP
		if err := s.Repo.InvalidateOTP(ctx, emailHash); err != nil {
			switch {
			// TODO: add all other possible cases
			default:
				return rp.NewValidationErr(errors.New("max attempts exceeded"), "OTP")
			}
		}
	}

	// Validate OTP
	if value != data.Code {
		// Increment attempts
		data.Attempts++
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return domain.NewJSONMarshalErr(err)
		}
		if err := s.Repo.StoreOTP(ctx, emailHash, dataBytes); err != nil {
			return err
		}
		return rp.NewValidationErr(errors.New("provided OTP code does not match stored OTP code"), "OTP")
	}

	// invalidate OTP if successful
	if err := s.Repo.InvalidateOTP(ctx, emailHash); err != nil {
		switch {
		// TODO: add all other possible cases
		default:
			return rp.NewValidationErr(errors.New("max attempts exceeded"), "OTP")
		}
	}
	return nil
}

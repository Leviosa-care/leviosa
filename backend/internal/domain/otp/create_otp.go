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

func (s *service) CreateOTP(ctx context.Context, emailHash string) (*OTP, error) {
	// get existing OTP
	otpEncoded, err := s.Repo.GetOTP(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return nil, err
		case errors.Is(err, rp.ErrNotFound):
			newOTP, err := s.newOTP(emailHash)
			if err != nil {
				return nil, domain.NewNotCreatedErr(err)
			}
			if err := s.crypto.ProcessStruct(ctx, newOTP); err != nil {
				return nil, domain.NewNotEncryptedErr("OTP", err)
			}
			otpData, err := json.Marshal(newOTP.Data.ToOTPEncrypted())
			if err != nil {
				return nil, domain.NewJSONMarshalErr(err)
			}
			if err := s.Repo.StoreOTP(ctx, emailHash, otpData); err != nil {
				switch {
				case errors.Is(err, rp.ErrDatabase):
					return nil, domain.NewQueryFailedErr(err)
				case errors.Is(err, rp.ErrContext):
					return nil, err
				}
			}
			return newOTP, nil
		}
	}
	var existingOTP OTP
	if err := json.Unmarshal(otpEncoded, &existingOTP); err != nil {
		return nil, domain.NewJSONUnmarshalErr(err)
	}

	// check if that otp is not expired
	if existingOTP.Data.Attempts != 0 && existingOTP.Data.Attempts < MaxOTPAttempts && time.Since(existingOTP.Data.CreatedAt) < time.Minute {
		return nil, domain.NewRateLimitErr(
			fmt.Errorf("please wait before requesting another OTP"),
			"otp",
		)
	}

	existingOTP.Data.Attempts++
	existingOTP.Data.CreatedAt = time.Now()

	otpData, err := json.Marshal(existingOTP)
	if err != nil {
		return nil, domain.NewJSONMarshalErr(err)
	}

	if err := s.Repo.StoreOTP(ctx, emailHash, otpData); err != nil {
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrContext):
			return nil, err
		}
	}
	return &existingOTP, nil
}

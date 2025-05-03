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
			newOTP, err := NewOTP(emailHash)
			if err != nil {
				return nil, domain.NewNotCreatedErr(err)
			}
			otpData, err := json.Marshal(newOTP)
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
	if existingOTP.Attempts != 0 && existingOTP.Attempts < MaxOTPAttempts && time.Since(existingOTP.Created) < time.Minute {
		return nil, domain.NewRateLimitErr(
			fmt.Errorf("please wait before requesting another OTP"),
			"otp",
		)
	}

	existingOTP.Attempts++
	existingOTP.Created = time.Now()

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

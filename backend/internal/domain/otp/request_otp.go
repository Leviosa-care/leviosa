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

func (s *service) RequestOTP(ctx context.Context, email string) (string, error) {
	emailHash := s.crypto.HashBasic(ctx, []byte(email))
	data, err := s.repo.GetOTP(ctx, emailHash)
	if err != nil && !errors.Is(err, rp.ErrNotFound) {
		return "", domain.NewQueryFailedErr(err)
	}
	if err == nil {
		var encrypted Data
		if err := json.Unmarshal(data, &encrypted); err != nil {
			return "", domain.NewJSONUnmarshalErr(err)
		}
		otp := &OTP{
			EmailHash:     emailHash,
			CodeEncrypted: encrypted.CodeEncrypted,
			Attempts:      encrypted.Attempts,
			ExpiresAt:     encrypted.ExpiresAt,
			CreatedAt:     encrypted.CreatedAt,
			DEKEncrypted:  encrypted.DEKEncrypted,
			KeyVersion:    encrypted.KeyVersion,
		}
		if err := s.crypto.DecryptStruct(ctx, otp); err != nil {
			return "", domain.NewNotDecryptedErr("OTP", err)
		}
		if !isOTPExpired(otp) && encrypted.Attempts < MaxOTPAttempts {
			return "", domain.NewRateLimitErr(fmt.Errorf("OTP recently requested"), "otp")
		}
	}
	// no existing OTP or expired - geneate new one
	otp, err := s.newOTP(email)
	if err != nil {
		return "", domain.NewNotCreatedErr(err)
	}
	if err := s.crypto.ProcessStruct(ctx, otp); err != nil {
		return "", domain.NewNotEncryptedErr("OTP", err)
	}
	encoded, err := json.Marshal(otp.Data())
	if err != nil {
		return "", domain.NewJSONMarshalErr(err)
	}
	if err := s.repo.SaveOTP(ctx, otp.EmailHash, encoded); err != nil {
		switch {
		case errors.Is(err, rp.ErrNotCreated):
			return "", domain.NewQueryFailedErr(err)
		default:
			return "", domain.NewNotCreatedErr(err)
		}
	}

	return otp.Code, nil
}

func isOTPExpired(otp *OTP) bool {
	return time.Since(otp.CreatedAt) > OTPDURATION
}

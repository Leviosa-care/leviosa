package otpRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *OTPRepository) SaveOTP(ctx context.Context, emailHash string, otpEncoded []byte, ttl time.Duration) error {
	key := FormatOTPKey(emailHash)

	if err := r.client.Set(ctx, key, otpEncoded, ttl).Err(); err != nil {
		return errs.ClassifyRedisError("save OTP", err)
	}

	return nil

}

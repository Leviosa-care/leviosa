package otpRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *OTPRepository) TouchOTP(ctx context.Context, emailHash string, ttl time.Duration) error {
	key := FormatOTPKey(emailHash)

	err := r.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return errs.ClassifyRedisError("update OTP TTL", err)
	}

	return nil
}

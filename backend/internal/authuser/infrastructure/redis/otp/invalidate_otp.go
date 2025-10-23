package otpRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *OTPRepository) InvalidateOTP(ctx context.Context, emailHash string) error {
	key := FormatOTPKey(emailHash)

	result := r.client.Del(ctx, key)
	if err := result.Err(); err != nil {
		return errs.ClassifyRedisError("delete OTP", err)
	}

	// Check if the key was actually deleted
	keysDeleted := result.Val()
	if keysDeleted == 0 {
		// Key didn't exist - another request already consumed the OTP
		return errs.ErrRepositoryNotFound
	}

	return nil
}

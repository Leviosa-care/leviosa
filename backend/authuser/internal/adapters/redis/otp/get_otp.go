package otpRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
)

func (r *OTPRepository) GetOTP(ctx context.Context, emailHash string) ([]byte, error) {
	key := FormatOTPKey(emailHash)

	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errs.ClassifyRedisError("get OTP", err)
	}

	return result, nil
}

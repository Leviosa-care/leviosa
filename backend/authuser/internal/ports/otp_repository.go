package ports

import (
	"context"
	"time"
)

type OTPRepository interface {
	SaveOTP(ctx context.Context, emailHash string, otpEncoded []byte, ttl time.Duration) error
	GetOTP(ctx context.Context, emailHash string) ([]byte, error)
	InvalidateOTP(ctx context.Context, emailHash string) error
}

package otpService

import (
	"context"
	"time"
)

type Reader interface {
	GetOTP(ctx context.Context, emailHash string) ([]byte, error)
}
type Writer interface {
	SaveOTP(ctx context.Context, emailHash string, otpEncoded []byte, duration time.Duration) error
	InvalidateOTP(ctx context.Context, emailHash string) error
	TouchOTP(ctx context.Context, emailHash string) error
}

type ReadWriter interface {
	Reader
	Writer
}

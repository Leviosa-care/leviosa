package otpService

import (
	"context"
)

type Reader interface {
	GetOTP(ctx context.Context, emailHash string) ([]byte, error)
}
type Writer interface {
	SaveOTP(ctx context.Context, emailHash string, otpEncoded []byte) error
	InvalidateOTP(ctx context.Context, emailHash string) error
}

type ReadWriter interface {
	Reader
	Writer
}

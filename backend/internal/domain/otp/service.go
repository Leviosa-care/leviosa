package otpService

import (
	"context"

	"github.com/hengadev/encx"
)

type Service interface {
	RequestOTP(ctx context.Context, email string) (string, error)
	VerifyOTP(ctx context.Context, email string, code string) error
	CancelOTP(ctx context.Context, email string) error
	ResendOTP(ctx context.Context, email string) (*OTP, error)
}

type service struct {
	repo   ReadWriter
	crypto encx.CryptoService
}

func New(repo ReadWriter, crypto encx.CryptoService) Service {
	return &service{repo: repo,
		crypto: crypto,
	}
}

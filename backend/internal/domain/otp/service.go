package otpService

import (
	"context"

	"github.com/hengadev/encx"
)

type Service interface {
	CancelOTP(ctx context.Context, email string) error
	CreateOTP(ctx context.Context, emailHash string) (*OTP, error)
	ValidateOTP(ctx context.Context, emailHash string, value string) error
}

type service struct {
	Repo   ReadWriter
	crypto encx.CryptoService
}

func New(repo ReadWriter, crypto encx.CryptoService) Service {
	return &service{
		Repo:   repo,
		crypto: crypto,
	}
}

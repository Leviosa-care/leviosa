package otpService

import "context"

type Service interface {
	CancelOTP(ctx context.Context, email string) error
	CreateOTP(ctx context.Context, emailHash string) (*OTP, error)
	ValidateOTP(ctx context.Context, emailHash string, value string) error
}

type service struct {
	Repo ReadWriter
}

func New(repo ReadWriter) Service {
	return &service{
		Repo: repo,
	}
}

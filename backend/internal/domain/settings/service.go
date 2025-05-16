package settings

import (
	"context"

	"github.com/hengadev/encx"
)

type Service interface { // getters
	GetCompanyName(ctx context.Context) (string, error)
	GetCompanyEmail(ctx context.Context) (string, error)
	GetCompanyTelephone(ctx context.Context) (string, error) // encrypted
	GetCompanyLegalAddress(ctx context.Context) (string, error)
	GetCompanyInstagram(ctx context.Context) (string, error)
	GetCompanyLogo(ctx context.Context) ([]byte, error)
	GetOTPDuration(ctx context.Context) (int, error)
	GetOTPLength(ctx context.Context) (int, error)
	GetOTPMaxAttempts(ctx context.Context) (int, error)
	// setters
	SetCompanyName(ctx context.Context, name string) error
	SetCompanyEmail(ctx context.Context, email string) error
	SetCompanyTelephone(ctx context.Context, telephone string) error // encrypted
	SetCompanyLegalAddress(ctx context.Context, address string) error
	SetCompanyInstagram(ctx context.Context, link string) error
	SetCompanyLogo(ctx context.Context, logo []byte) error
	SetOTPDuration(ctx context.Context, duration int) error
	SetOTPLength(ctx context.Context, length int) error
	SetOTPMaxAttempts(ctx context.Context, value int) error
}

type service struct {
	repo         readWriter
	media        mediaReadWriter
	crypto       encx.CryptoService
}

func New(
	repo readWriter,
	media mediaReadWriter,
	crypto encx.CryptoService,
) Service {
	return &service{
		repo:   repo,
		media:  media,
		crypto: crypto,
	}
}

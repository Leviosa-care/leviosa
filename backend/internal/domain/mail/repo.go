package mailService

import (
	"context"
)

type Reader interface {
	GetLogo(ctx context.Context) (string, error)
	GetCompanyEmail(ctx context.Context) (string, error)
	GetCompanyLegalAddress(ctx context.Context) (string, error)
	GetCompanyInstagram(ctx context.Context) (string, error)
}

package ports

import (
	"context"
)

type SettingsRepository interface {
	GetCompanyEmail(ctx context.Context) (string, error)
	GetCompanyLegalAddress(ctx context.Context) (string, error)
	GetCompanyInstagram(ctx context.Context) (string, error)
	GetCompanyLogo(ctx context.Context) ([]byte, error)
}


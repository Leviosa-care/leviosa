package ports

import "context"

// SettingsProvider provides access to company settings with caching
// This port allows the notification service to access settings without direct coupling
type SettingsProvider interface {
	GetCompanyEmail(ctx context.Context) (string, error)
	GetCompanyLegalAddress(ctx context.Context) (string, error)
	GetCompanyInstagram(ctx context.Context) (string, error)
	GetCompanyLogo(ctx context.Context) ([]byte, error)
	GetCompanyTelephone(ctx context.Context) (string, error)
	GetCompanyName(ctx context.Context) (string, error)

	// Cache management for RabbitMQ consumer
	InvalidateCache(key string)
	InvalidateAllCache()
}

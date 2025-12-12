package rabbitmq

import (
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

// SettingsConsumer listens to settings updates via RabbitMQ and warms the cache
type SettingsConsumer struct {
	settingsProvider ports.SettingsProvider
}

func NewSettingsConsumer(settingsProvider ports.SettingsProvider) *SettingsConsumer {
	return &SettingsConsumer{
		settingsProvider: settingsProvider,
	}
}

// ConsumeSettingsUpdate processes settings update messages
func (c *SettingsConsumer) ConsumeSettingsUpdate(dto settings.SettingDTO) error {
	// Invalidate cache key to force fresh fetch on next access
	// This ensures data consistency between settings service and notification cache

	switch dto.Key {
	case settings.CompanyEmail:
		c.settingsProvider.InvalidateCache("company_email")

	case settings.CompanyLegalAddress:
		c.settingsProvider.InvalidateCache("company_address")

	case settings.CompanyInstagram:
		c.settingsProvider.InvalidateCache("company_instagram")

	case settings.CompanyLogo:
		c.settingsProvider.InvalidateCache("company_logo")

	case settings.CompanyPhone:
		c.settingsProvider.InvalidateCache("company_telephone")

	case settings.CompanyName:
		c.settingsProvider.InvalidateCache("company_name")

	default:
		// Ignore settings that notification service doesn't use
		return nil
	}

	return nil
}

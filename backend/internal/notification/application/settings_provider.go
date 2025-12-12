package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
	settingsPorts "github.com/Leviosa-care/leviosa/backend/internal/settings/ports"
)

// settingsProvider implements ports.SettingsProvider
type settingsProvider struct {
	settingsService settingsPorts.SettingsService
	cache           *domain.CompanyCache
	cacheTTL        time.Duration
}

// NewSettingsProvider creates a new settings provider with caching
func NewSettingsProvider(settingsService settingsPorts.SettingsService) ports.SettingsProvider {
	return &settingsProvider{
		settingsService: settingsService,
		cache:           domain.NewCompanyCache(),
		cacheTTL:        5 * time.Minute,
	}
}

func (sp *settingsProvider) GetCompanyEmail(ctx context.Context) (string, error) {
	if cached, valid := sp.cache.Get("company_email"); valid {
		return cached.(string), nil
	}

	resp, err := sp.settingsService.GetCompanyEmail(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get company email: %w", err)
	}

	sp.cache.Set("company_email", resp.Email, sp.cacheTTL)
	return resp.Email, nil
}

func (sp *settingsProvider) GetCompanyLegalAddress(ctx context.Context) (string, error) {
	if cached, valid := sp.cache.Get("company_address"); valid {
		return cached.(string), nil
	}

	resp, err := sp.settingsService.GetCompanyLegalAddress(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get company address: %w", err)
	}

	sp.cache.Set("company_address", resp.Address, sp.cacheTTL)
	return resp.Address, nil
}

func (sp *settingsProvider) GetCompanyInstagram(ctx context.Context) (string, error) {
	if cached, valid := sp.cache.Get("company_instagram"); valid {
		return cached.(string), nil
	}

	resp, err := sp.settingsService.GetCompanyInstagram(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get company instagram: %w", err)
	}

	sp.cache.Set("company_instagram", resp.Instagram, sp.cacheTTL)
	return resp.Instagram, nil
}

func (sp *settingsProvider) GetCompanyLogo(ctx context.Context) ([]byte, error) {
	// TODO: Settings service returns LogoURL not LogoData
	// Need to implement fetching logo from S3/URL or change interface to return URL
	// For now, return empty to allow compilation
	return nil, nil
}

func (sp *settingsProvider) GetCompanyTelephone(ctx context.Context) (string, error) {
	if cached, valid := sp.cache.Get("company_telephone"); valid {
		return cached.(string), nil
	}

	resp, err := sp.settingsService.GetCompanyTelephone(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get company telephone: %w", err)
	}

	sp.cache.Set("company_telephone", resp.Telephone, sp.cacheTTL)
	return resp.Telephone, nil
}

func (sp *settingsProvider) GetCompanyName(ctx context.Context) (string, error) {
	if cached, valid := sp.cache.Get("company_name"); valid {
		return cached.(string), nil
	}

	resp, err := sp.settingsService.GetCompanyName(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get company name: %w", err)
	}

	sp.cache.Set("company_name", resp.Name, sp.cacheTTL)
	return resp.Name, nil
}

// InvalidateCache removes a specific key from cache (used by RabbitMQ consumer)
func (sp *settingsProvider) InvalidateCache(key string) {
	sp.cache.Invalidate(key)
}

// InvalidateAllCache clears entire cache
func (sp *settingsProvider) InvalidateAllCache() {
	sp.cache.InvalidateAll()
}

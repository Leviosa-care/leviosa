package oauth

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/nextcloud"
)

// Config holds OAuth configuration
type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	AppleClientID      string
	AppleClientSecret  string
	AppleTeamID        string
	AppleKeyID         string
	ApplePrivateKey    string
	SessionSecret      string
	BaseURL            string
	// Test-only fields
	UseNextcloudForTesting bool
	NextcloudTestURL       string
}

// InitializeOAuthProviders initializes OAuth providers using environment variables
func InitializeOAuthProviders() error {
	config := LoadConfigFromEnv()
	return InitializeProviders(config)
}

// InitializeProviders sets up OAuth providers using Goth
func InitializeProviders(config *Config) error {
	// Initialize session store
	key := config.SessionSecret
	if key == "" {
		return fmt.Errorf("SESSION_SECRET environment variable is required")
	}

	maxAge := 86400 * 30 // 30 days
	isProd := os.Getenv("ENV") == "production"

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	var providers []goth.Provider

	// Configure Google provider (or Nextcloud for testing)
	if config.GoogleClientID != "" && config.GoogleClientSecret != "" {
		if config.UseNextcloudForTesting && config.NextcloudTestURL != "" {
			// For testing: use Nextcloud but register as "google" provider
			googleProvider := nextcloud.NewCustomisedDNS(
				config.GoogleClientID,
				config.GoogleClientSecret,
				config.BaseURL+"/auth/oauth/google/callback",
				config.NextcloudTestURL,
			)
			// Override the provider name to be "google" 
			providers = append(providers, renameProvider(googleProvider, "google"))
		} else {
			// Production: use real Google provider
			googleProvider := google.New(
				config.GoogleClientID,
				config.GoogleClientSecret,
				config.BaseURL+"/auth/oauth/google/callback",
				"email", "profile",
			)
			providers = append(providers, googleProvider)
		}
	}

	// Configure Apple provider
	if config.AppleClientID != "" && config.AppleClientSecret != "" {
		appleProvider := apple.New(
			config.AppleClientID,
			config.AppleClientSecret,
			config.BaseURL+"/auth/oauth/apple/callback",
			nil, // scopes - Apple doesn't use traditional scopes
			config.AppleTeamID,
			config.AppleKeyID,
			config.ApplePrivateKey,
		)
		providers = append(providers, appleProvider)
	}


	if len(providers) == 0 {
		return fmt.Errorf("no OAuth providers configured")
	}

	goth.UseProviders(providers...)
	return nil
}

// LoadConfigFromEnv loads OAuth configuration from environment variables
func LoadConfigFromEnv() *Config {
	return &Config{
		GoogleClientID:         os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:     os.Getenv("GOOGLE_CLIENT_SECRET"),
		AppleClientID:          os.Getenv("APPLE_CLIENT_ID"),
		AppleClientSecret:      os.Getenv("APPLE_CLIENT_SECRET"),
		AppleTeamID:            os.Getenv("APPLE_TEAM_ID"),
		AppleKeyID:             os.Getenv("APPLE_KEY_ID"),
		ApplePrivateKey:        os.Getenv("APPLE_PRIVATE_KEY"),
		SessionSecret:          os.Getenv("SESSION_SECRET"),
		BaseURL:                getBaseURL(),
		UseNextcloudForTesting: os.Getenv("USE_NEXTCLOUD_FOR_TESTING") == "true",
		NextcloudTestURL:       os.Getenv("NEXTCLOUD_TEST_URL"),
	}
}

func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5000" // development default
	}
	return baseURL
}

// renameProvider wraps a provider to override its name
// This allows us to register Nextcloud as "google" for testing
func renameProvider(provider goth.Provider, newName string) goth.Provider {
	return &providerWrapper{
		Provider: provider,
		name:     newName,
	}
}

// providerWrapper wraps a provider to override its name
type providerWrapper struct {
	goth.Provider
	name string
}

// Name returns the overridden name
func (p *providerWrapper) Name() string {
	return p.name
}


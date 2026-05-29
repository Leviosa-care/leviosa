package app

import (
	"context"
	"testing"
)

func TestLoadConfig_RejectsEmptySessionSecretInStaging(t *testing.T) {
	t.Setenv("ENVIRONMENT", "staging")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	// SESSION_SECRET intentionally not set → defaults to "development-secret-key"

	_, err := LoadConfig(context.Background())
	if err == nil {
		t.Fatal("expected LoadConfig to reject empty/default SESSION_SECRET in staging, got nil error")
	}
}

func TestLoadConfig_RejectsDefaultSessionSecretInProduction(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("SESSION_SECRET", "development-secret-key")

	_, err := LoadConfig(context.Background())
	if err == nil {
		t.Fatal("expected LoadConfig to reject default SESSION_SECRET in production, got nil error")
	}
}

func TestLoadConfig_AcceptsRealSecretInStaging(t *testing.T) {
	t.Setenv("ENVIRONMENT", "staging")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("SESSION_SECRET", "a-real-secret-that-is-not-the-default")

	cfg, err := LoadConfig(context.Background())
	if err != nil {
		t.Fatalf("expected LoadConfig to succeed with a real secret in staging, got error: %v", err)
	}
	if cfg.SessionSecret != "a-real-secret-that-is-not-the-default" {
		t.Errorf("expected SessionSecret to be set, got %q", cfg.SessionSecret)
	}
}

func TestLoadConfig_DevelopmentAllowsDefaultSecret(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	// Ensure SESSION_SECRET is not set in the test environment
	t.Setenv("SESSION_SECRET", "")

	cfg, err := LoadConfig(context.Background())
	if err != nil {
		t.Fatalf("expected LoadConfig to succeed in development with default secret, got error: %v", err)
	}
	if cfg.SessionSecret != "development-secret-key" {
		t.Errorf("expected default SessionSecret in development, got %q", cfg.SessionSecret)
	}
}

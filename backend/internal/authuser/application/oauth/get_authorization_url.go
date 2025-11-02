package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/markbates/goth"
)

// GetAuthorizationURL generates an authorization URL for the given provider
func (s *Service) GetAuthorizationURL(ctx context.Context, provider string) (string, string, error) {
	// Get the provider from Goth
	gothProvider, err := goth.GetProvider(provider)
	if err != nil {
		return "", "", fmt.Errorf("provider %s not found: %w", provider, err)
	}

	// Generate a random state parameter for CSRF protection
	state, err := generateState()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Create a session to get the authorization URL
	sess, err := gothProvider.BeginAuth(state)
	if err != nil {
		return "", "", fmt.Errorf("failed to begin auth: %w", err)
	}

	authURL, err := sess.GetAuthURL()
	if err != nil {
		return "", "", fmt.Errorf("failed to get auth URL: %w", err)
	}

	return authURL, state, nil
}

// generateState generates a random state parameter for CSRF protection
func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// Service handles OAuth operations
type Service struct{}

// NewService creates a new OAuth service
func NewService() *Service {
	return &Service{}
}

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

// ExchangeCodeForUser exchanges authorization code for user information using HTTP request
func (s *Service) ExchangeCodeForUser(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	// Use gothic's built-in method to complete the authentication
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return goth.User{}, fmt.Errorf("failed to complete user auth: %w", err)
	}

	return user, nil
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


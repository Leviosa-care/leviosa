package oauth

import (
	"context"
	"fmt"

	"github.com/markbates/goth"
)

// GetAuthorizationURLWithState generates an authorization URL using a caller-supplied state
// instead of generating one internally.  Use this when the caller needs to embed
// additional context (e.g. a user ID for the link flow) into the state.
func (s *Service) GetAuthorizationURLWithState(_ context.Context, provider, state string) (string, error) {
	gothProvider, err := goth.GetProvider(provider)
	if err != nil {
		return "", fmt.Errorf("provider %s not found: %w", provider, err)
	}

	sess, err := gothProvider.BeginAuth(state)
	if err != nil {
		return "", fmt.Errorf("failed to begin auth: %w", err)
	}

	authURL, err := sess.GetAuthURL()
	if err != nil {
		return "", fmt.Errorf("failed to get auth URL: %w", err)
	}

	return authURL, nil
}

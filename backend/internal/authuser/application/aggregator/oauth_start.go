package aggregator

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/oauth"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// OAuthStart initiates the OAuth flow and returns the authorization URL
func (s *AuthAggregatorService) OAuthStart(ctx context.Context, request *domain.OAuthStartRequest) (*domain.OAuthStartResponse, error) {
	// Validate the request
	if request.Provider == "" {
		return nil, errs.NewInvalidValueErr("provider is required")
	}

	// Validate provider is supported
	if request.Provider != "google" && request.Provider != "apple" {
		return nil, errs.NewInvalidValueErr("unsupported provider: " + request.Provider)
	}

	// Create OAuth service instance
	oauthService := oauth.NewService()

	// Get authorization URL
	authURL, state, err := oauthService.GetAuthorizationURL(ctx, request.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization URL: %w", err)
	}

	return &domain.OAuthStartResponse{
		AuthorizationURL: authURL,
		State:            state,
	}, nil
}
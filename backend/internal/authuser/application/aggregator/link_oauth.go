package aggregator

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/oauth"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// LinkOAuth initiates the OAuth flow for linking a provider to an existing authenticated user.
// It returns an authorization URL that the frontend should redirect to.
func (s *AuthAggregatorService) LinkOAuth(ctx context.Context, sessionInfo *session.SessionInfo, provider string) (*domain.OAuthLinkResponse, error) {
	if provider != "google" && provider != "apple" {
		return nil, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	// Fetch the current user to check if already linked
	user, err := s.user.GetUserByID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user for OAuth link: %w", err)
	}

	// Check if this provider is already linked
	switch provider {
	case "google":
		if user.GoogleID != "" {
			return nil, errs.NewConflictErr(fmt.Errorf("%s account is already linked", provider))
		}
	case "apple":
		if user.AppleID != "" {
			return nil, errs.NewConflictErr(fmt.Errorf("%s account is already linked", provider))
		}
	}

	// Create OAuth service instance
	oauthService := oauth.NewService()

	// Get authorization URL
	authURL, state, err := oauthService.GetAuthorizationURL(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization URL for link: %w", err)
	}

	return &domain.OAuthLinkResponse{
		AuthorizationURL: authURL,
		State:            state,
	}, nil
}

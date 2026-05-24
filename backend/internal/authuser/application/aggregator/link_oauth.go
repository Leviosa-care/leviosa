package aggregator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/oauth"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// LinkOAuth initiates the OAuth flow for linking a provider to an existing authenticated user.
// The state parameter embeds the user ID so the callback handler can route to CompleteLinkOAuth
// instead of the normal sign-in path.
func (s *AuthAggregatorService) LinkOAuth(ctx context.Context, sessionInfo *session.SessionInfo, provider string) (*domain.OAuthLinkResponse, error) {
	if provider != "google" && provider != "apple" {
		return nil, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	user, err := s.user.GetUserByID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user for OAuth link: %w", err)
	}

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

	// Embed user ID in state so the callback can distinguish a link flow from a sign-in.
	// Format: "link.{userID}.{base64random}"
	random := make([]byte, 24)
	if _, err := rand.Read(random); err != nil {
		return nil, fmt.Errorf("generate link state: %w", err)
	}
	state := "link." + sessionInfo.UserID.String() + "." + base64.RawURLEncoding.EncodeToString(random)

	oauthService := oauth.NewService()
	authURL, err := oauthService.GetAuthorizationURLWithState(ctx, provider, state)
	if err != nil {
		return nil, fmt.Errorf("get authorization URL for link: %w", err)
	}

	return &domain.OAuthLinkResponse{
		AuthorizationURL: authURL,
		State:            state,
	}, nil
}

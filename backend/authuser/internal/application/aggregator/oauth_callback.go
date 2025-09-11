package aggregator

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/authuser/internal/application/oauth"
	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

// OAuthCallback handles the OAuth callback and creates/logs in the user
func (s *AuthAggregatorService) OAuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request, provider string) (*domain.OAuthCallbackResponse, error) {
	// Validate provider is supported
	if provider != "google" && provider != "apple" {
		return nil, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	// Create OAuth service instance
	oauthService := oauth.NewService()

	// Exchange authorization code for user information
	gothUser, err := oauthService.ExchangeCodeForUser(w, r)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for user: %w", err)
	}

	// Extract user information from OAuth response
	var oauthUserID string
	switch provider {
	case "google":
		oauthUserID = gothUser.UserID
	case "apple":
		oauthUserID = gothUser.UserID
	default:
		return nil, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	// Get or create user
	user, isNewUser, err := s.user.GetOrCreateOAuthUser(
		ctx,
		provider,
		oauthUserID,
		gothUser.Email,
		gothUser.FirstName,
		gothUser.LastName,
	)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidValue) {
			return nil, errs.NewInvalidValueErr("invalid OAuth user data")
		}
		return nil, fmt.Errorf("failed to get or create OAuth user: %w", err)
	}

	// Check if user is in active state (not pending or unverified)
	if user.State != domain.Active {
		return nil, errs.NewUnauthorizedErr("account not activated")
	}

	// Parse user ID
	userUUID, err := uuid.Parse(user.ID.String())
	if err != nil {
		return nil, errs.NewInternalErr(fmt.Errorf("invalid user ID format: %w", err))
	}

	// Convert role
	role, ok := identity.ConvertToRole(user.Role)
	if !ok {
		return nil, errs.NewInternalErr(fmt.Errorf("invalid user role: %s", user.Role))
	}

	// Create session
	token, err := s.session.CreateSession(ctx, &domain.CreateSessionRequest{
		UserID: userUUID.String(),
		Role:   role,
	})
	if err != nil {
		if errors.Is(err, errs.ErrInvalidValue) {
			return nil, errs.NewInvalidValueErr("invalid session data")
		}
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.OAuthCallbackResponse{
		AccessToken:        token.AccessToken,
		RefreshToken:       token.RefreshToken,
		AccessTokenExpiry:  token.AccessTokenExpiry.Unix(),
		RefreshTokenExpiry: token.RefreshTokenExpiry.Unix(),
		IsNewUser:          isNewUser,
	}, nil
}
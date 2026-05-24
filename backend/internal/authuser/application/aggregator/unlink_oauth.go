package aggregator

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// UnlinkOAuth removes the association between an OAuth provider and the authenticated user.
// It guards against unlinking the only sign-in method when the user has no password set.
func (s *AuthAggregatorService) UnlinkOAuth(ctx context.Context, sessionInfo *session.SessionInfo, provider string) (*domain.OAuthUnlinkResponse, error) {
	if provider != "google" && provider != "apple" {
		return nil, errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	// Fetch the current user
	user, err := s.user.GetUserByID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user for OAuth unlink: %w", err)
	}

	// Verify the provider is currently linked
	switch provider {
	case "google":
		if user.GoogleID == "" {
			return nil, errs.NewInvalidValueErr(provider + " account is not linked")
		}
	case "apple":
		if user.AppleID == "" {
			return nil, errs.NewInvalidValueErr(provider + " account is not linked")
		}
	}

	// Guard: cannot unlink the only sign-in method if no password is set.
	// Count remaining OAuth providers after unlinking this one.
	remainingProviders := 0
	if user.GoogleID != "" && provider != "google" {
		remainingProviders++
	}
	if user.AppleID != "" && provider != "apple" {
		remainingProviders++
	}

	// If no other OAuth providers remain, the user must have a password
	if remainingProviders == 0 && !user.HasPassword {
		return nil, errs.NewUnprocessableEntityErr("cannot unlink the only sign-in method; please set a password first")
	}

	// Clear the provider ID on the user record
	updateReq := &domain.UpdateUserRequest{}
	switch provider {
	case "google":
		empty := ""
		updateReq.GoogleID = &empty
	case "apple":
		empty := ""
		updateReq.AppleID = &empty
	}

	_, err = s.user.UpdateUser(ctx, sessionInfo.UserID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("unlink OAuth account: %w", err)
	}

	return &domain.OAuthUnlinkResponse{
		Provider: provider,
	}, nil
}

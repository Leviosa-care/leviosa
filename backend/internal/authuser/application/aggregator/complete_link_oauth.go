package aggregator

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/oauth"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// CompleteLinkOAuth finishes an OAuth link flow by associating the provider identity
// with the authenticated user.  Called from the OAuth callback when the state begins
// with "link." (see LinkOAuth for the state format).
func (s *AuthAggregatorService) CompleteLinkOAuth(ctx context.Context, userID uuid.UUID, provider string, w http.ResponseWriter, r *http.Request) error {
	if provider != "google" && provider != "apple" {
		return errs.NewInvalidValueErr("unsupported provider: " + provider)
	}

	oauthService := oauth.NewService()
	gothUser, err := oauthService.ExchangeCodeForUser(w, r)
	if err != nil {
		return fmt.Errorf("exchange code for user: %w", err)
	}

	providerUserID := gothUser.UserID

	// Ensure this provider identity is not already linked to a different account.
	switch provider {
	case "google":
		existing, err := s.user.GetUserByGoogleID(ctx, providerUserID)
		if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
			return fmt.Errorf("check existing Google link: %w", err)
		}
		if existing != nil && existing.ID != userID {
			return errs.NewConflictErr(fmt.Errorf("this Google account is already linked to another user"))
		}
	case "apple":
		existing, err := s.user.GetUserByAppleID(ctx, providerUserID)
		if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
			return fmt.Errorf("check existing Apple link: %w", err)
		}
		if existing != nil && existing.ID != userID {
			return errs.NewConflictErr(fmt.Errorf("this Apple account is already linked to another user"))
		}
	}

	updateReq := &domain.UpdateUserRequest{}
	switch provider {
	case "google":
		updateReq.GoogleID = &providerUserID
	case "apple":
		updateReq.AppleID = &providerUserID
	}

	if _, err := s.user.UpdateUser(ctx, userID, updateReq); err != nil {
		return fmt.Errorf("link OAuth account: %w", err)
	}

	return nil
}

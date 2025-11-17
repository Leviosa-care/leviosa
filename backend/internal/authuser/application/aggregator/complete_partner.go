package aggregator

import (
	"context"
	"errors"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// CompletePartner handles partner registration completion.
//
// This method combines user profile completion with partner profile creation in a single operation.
// It is called during the partner registration flow after email verification.
//
// Flow:
//  1. Verify session is in pending state
//  2. Complete user profile (name, password, address, etc.) via user service
//  3. Create partner profile (bio, categories, products) via partner service
//  4. Mark session as completed
//  5. Revoke all sessions to force re-login after admin approval
//
// The user and partner both start in unverified/pending state and require admin approval.
func (s *AuthAggregatorService) CompletePartner(ctx context.Context, sessionInfo *session.SessionInfo, request *domain.CompletePartnerRequest) error {
	// Verify session is in pending state
	if sessionInfo.State != session.SessionPending {
		return errs.NewConflictErr(errors.New("session is not in pending state"))
	}

	// Complete the user information (creates Stripe customer, sets user data, State=Pending)
	userRequest := request.ToCompleteUserRequest()
	if err := s.user.CompleteUser(ctx, sessionInfo.UserID, userRequest); err != nil {
		return err
	}

	// Create partner profile (IsVerified=false, requires admin approval)
	if _, err := s.partner.CreatePartner(
		ctx,
		sessionInfo.UserID,
		request.Bio,
		request.Experience,
		// request.Certifications,
		request.CategoryIDs,
		request.ProductIDs,
	); err != nil {
		return err
	}

	// Partner completed successfully - mark completion timestamp in session
	completedAt := time.Now()
	if err := s.session.UpdateSessionCompletion(ctx, sessionInfo.ID, &completedAt); err != nil {
		return err
	}

	// Remove sessions to force re-authentication after admin approval
	if err := s.session.RevokeAllUserSessions(ctx, sessionInfo.UserID); err != nil {
		return err
	}

	return nil
}

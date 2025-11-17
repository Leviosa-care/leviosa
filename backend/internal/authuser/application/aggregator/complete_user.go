package aggregator

import (
	"context"
	"errors"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) CompleteUser(ctx context.Context, sessionInfo *session.SessionInfo, request *domain.CompleteUserRequest) error {
	// Verify session is in pending state
	if sessionInfo.State != session.SessionPending {
		return errs.NewConflictErr(errors.New("session is not in pending state"))
	}

	// Complete the user information
	if err := s.user.CompleteUser(ctx, sessionInfo.UserID, request); err != nil {
		return err
	}

	// User completed successfully - mark completion timestamp in session
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

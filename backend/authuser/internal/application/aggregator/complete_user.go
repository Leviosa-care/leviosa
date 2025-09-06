package aggregator

import (
	"context"
	"errors"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
)

func (s *AuthAggregatorService) CompleteUser(ctx context.Context, sessionInfo *session.SessionInfo, request *domain.CompleteUserRequest) error {
	// Verify session is in pending state
	if sessionInfo.State != session.SessionPending {
		return errs.NewConflictErr(errors.New("session is not in pending state"))
	}

	// Complete the user information
	if err := s.user.CompleteUser(ctx, sessionInfo.UserID, request); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return err // Pass through validation errors
		case errors.Is(err, errs.ErrDomainNotFound):
			return err // Pass through not found errors (user doesn't exist)
		case errors.Is(err, errs.ErrConflict):
			return err // Pass through conflict errors (user already completed or wrong state)
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors (database issues)
		default:
			return errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	// User completed successfully - mark completion timestamp in session
	completedAt := time.Now()
	if err := s.session.UpdateSessionCompletion(ctx, sessionInfo.ID, &completedAt); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return err // Pass through validation errors
		case errors.Is(err, errs.ErrDomainNotFound):
			return err // Pass through not found errors
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		default:
			return errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	// Remove sessions to force re-authentication after admin approval
	if err := s.session.RevokeAllUserSessions(ctx, sessionInfo.UserID); err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			// Session already removed - this is acceptable
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		default:
			return errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	return nil
}

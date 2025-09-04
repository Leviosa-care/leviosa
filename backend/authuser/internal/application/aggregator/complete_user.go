package aggregator

import (
	"context"
	"errors"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (s *AuthAggregatorService) CompleteUser(ctx context.Context, sessionToken string, request *domain.CompleteUserRequest) error {
	// Get the current session using the token
	getSessionRequest := &domain.GetSessionRequest{Token: sessionToken}
	session, err := s.session.GetSession(ctx, getSessionRequest)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return err // Pass through validation errors
		case errors.Is(err, errs.ErrDomainNotFound):
			return err // Pass through not found errors (session doesn't exist)
		case errors.Is(err, errs.ErrExpiredToken):
			return err // Pass through expired token errors
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors (database issues)
		default:
			return errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	// Verify session is in pending state
	if session.State != auth.SessionPending {
		return errs.NewConflictErr(errors.New("session is not in pending state"))
	}

	// Complete the user information
	if err := s.user.CompleteUser(ctx, session.UserID, request); err != nil {
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
	if err := s.session.UpdateSessionCompletion(ctx, session.ID.String(), &completedAt); err != nil {
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

	// Remove session to force re-authentication after admin approval
	removeSessionRequest := &domain.RemoveSessionRequest{Token: sessionToken}
	if err := s.session.RemoveSession(ctx, removeSessionRequest); err != nil {
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

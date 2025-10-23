package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) DeleteOwnAccount(ctx context.Context, sessionInfo *session.SessionInfo) error {
	// Get user details for cleanup operations (email for OTP cleanup)
	userResponse, err := s.user.GetUserByID(ctx, sessionInfo.UserID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			return err // Pass through - user doesn't exist
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		default:
			return errs.NewUnexpectedError(err) // Wrap unexpected errors
		}
	}

	// 1. Revoke all user sessions (Redis cleanup) - including current session
	if err := s.session.RevokeAllUserSessions(ctx, sessionInfo.UserID); err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			// No sessions to revoke - this is acceptable
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		default:
			return errs.NewUnexpectedError(err) // Wrap unexpected errors
		}
	}

	// 2. Cancel/invalidate any pending OTPs for this user
	if err := s.otp.CancelOTP(ctx, userResponse.Email); err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			// No OTP to cancel - this is acceptable
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		default:
			return errs.NewUnexpectedError(err) // Wrap unexpected errors
		}
	}

	// 3. Delete the user record (this coordinates Stripe customer deletion automatically)
	if err := s.user.DeleteUser(ctx, sessionInfo.UserID); err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			return err // Pass through - user doesn't exist
		case errors.Is(err, errs.ErrInvalidValue):
			return err // Pass through validation errors
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		default:
			return errs.NewUnexpectedError(err) // Wrap unexpected errors
		}
	}

	return nil
}

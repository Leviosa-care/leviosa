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
		return err
	}

	// 1. Revoke all user sessions (Redis cleanup) - including current session
	if err := s.session.RevokeAllUserSessions(ctx, sessionInfo.UserID); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// No sessions to revoke - this is acceptable
		default:
			return err
		}
	}

	// 2. Cancel/invalidate any pending OTPs for this user
	if err := s.otp.CancelOTP(ctx, userResponse.Email); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// No OTP to cancel - this is acceptable
		default:
			return err
		}
	}

	// 3. Delete the user record (this coordinates Stripe customer deletion automatically)
	if err := s.user.DeleteUser(ctx, sessionInfo.UserID); err != nil {
		return err
	}

	return nil
}

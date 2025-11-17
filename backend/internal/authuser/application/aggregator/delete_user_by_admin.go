package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *AuthAggregatorService) DeleteUserByAdmin(ctx context.Context, userID uuid.UUID) error {
	// Get user first to retrieve email for OTP cleanup
	userResponse, err := s.user.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 1. Revoke all user sessions (Redis cleanup)
	if err := s.session.RevokeAllUserSessions(ctx, userID); err != nil {
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
	if err := s.user.DeleteUser(ctx, userID); err != nil {
		return err
	}

	return nil
}

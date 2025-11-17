package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) ConfirmPasswordReset(ctx context.Context, request *domain.ConfirmPasswordResetRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Validate and consume reset session (single-use)
	userEmail, err := s.session.ValidateResetSession(ctx, request.Token)
	if err != nil {
		return err
	}

	// Get user by email to update password
	user, err := s.user.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	// Reset the user's password (bypasses old password verification)
	if err := s.user.ResetPassword(ctx, user.ID, request.NewPassword); err != nil {
		return err
	}

	return nil
}

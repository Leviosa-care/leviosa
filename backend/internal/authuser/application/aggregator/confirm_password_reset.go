package aggregator

import (
	"context"
	"errors"
	"fmt"

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
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			return err // Pass through not found errors (invalid or expired token)
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections), errors.Is(err, errs.ErrResourceExhausted), errors.Is(err, errs.ErrTransactionFailure):
			return err // Pass through infrastructure errors
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("validate reset session cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("validate reset session timeout: %w", err)
		default:
			return fmt.Errorf("validate reset session: %w", err) // Wrap unexpected errors with context
		}
	}

	// Get user by email to update password
	user, err := s.user.GetUserByEmail(ctx, userEmail)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			// User not found - this shouldn't happen if reset session was valid
			return fmt.Errorf("user not found during password reset: %w", err)
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections), errors.Is(err, errs.ErrResourceExhausted), errors.Is(err, errs.ErrTransactionFailure):
			return err // Pass through infrastructure errors
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("get user for password reset cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("get user for password reset timeout: %w", err)
		default:
			return fmt.Errorf("get user for password reset: %w", err) // Wrap unexpected errors with context
		}
	}

	// Reset the user's password (bypasses old password verification)
	if err := s.user.ResetPassword(ctx, user.ID, request.NewPassword); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return err // Pass through validation errors (invalid new password)
		case errors.Is(err, errs.ErrDomainNotFound):
			return fmt.Errorf("user not found during password reset: %w", err)
		case errors.Is(err, errs.ErrExternalService):
			return err // Pass through external service errors
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections), errors.Is(err, errs.ErrResourceExhausted), errors.Is(err, errs.ErrTransactionFailure):
			return err // Pass through infrastructure errors
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("reset password cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("reset password timeout: %w", err)
		default:
			return fmt.Errorf("reset password: %w", err) // Wrap unexpected errors with context
		}
	}

	return nil
}

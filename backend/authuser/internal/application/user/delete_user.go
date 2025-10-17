package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// 1. Get user first to retrieve encrypted Stripe customer ID
	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to get user: %w", err))
		}
	}

	// 2. Decrypt user data to access Stripe customer ID
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("user", err)
	}

	// 3. Delete Stripe customer if exists
	if user.StripeCustomerID != "" {
		_, err := s.stripe.DeleteCustomer(ctx, user.StripeCustomerID)
		if err != nil {
			switch {
			case errors.Is(err, errs.ErrInvalidValue):
				// Stripe customer not found or invalid - this is acceptable
				// The customer might have been deleted already or never existed
			case errors.Is(err, errs.ErrPermissionDenied):
				return errs.NewPermissionErr(fmt.Sprintf("stripe customer deletion failed: %s", err.Error()))
			case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrResourceExhausted):
				return errs.NewExternalServiceErr(err, "stripe service unavailable")
			case errors.Is(err, errs.ErrRateLimit):
				return errs.NewExternalServiceErr(err, "stripe rate limit exceeded")
			default:
				return errs.NewUnexpectedError(fmt.Errorf("failed to delete stripe customer: %w", err))
			}
		}
	}

	// 4. Delete user from database
	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to delete user: %w", err))
		}
	}

	return nil
}


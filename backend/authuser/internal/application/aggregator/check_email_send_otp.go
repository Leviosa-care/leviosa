package aggregator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (s *AuthAggregatorService) CheckEmailSendOTP(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	available, err := s.user.CheckEmailAvailability(ctx, request)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// Email validation failed
			return fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections), errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock), errors.Is(err, errs.ErrResourceExhausted), errors.Is(err, errs.ErrPermissionDenied), errors.Is(err, errs.ErrDatabase):
			// Database infrastructure errors
			return fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return fmt.Errorf("check email availability cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return fmt.Errorf("check email availability timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return fmt.Errorf("check email availability: %w", err)
		}
	}

	if !available {
		return errs.NewConflictErr(errors.New("email is already registered"))
	}

	if err := s.otp.RequestOTP(ctx, request.Email); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue), errors.Is(err, errs.ErrRateLimit):
			// Invalid parameters or rate limiting
			return fmt.Errorf("request OTP: %w", err)
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections), errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrResourceExhausted), errors.Is(err, errs.ErrExternalService):
			// Infrastructure and external service errors
			return fmt.Errorf("request OTP: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return fmt.Errorf("request OTP cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return fmt.Errorf("request OTP timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return fmt.Errorf("request OTP: %w", err)
		}
	}
	return nil
}

package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (s *UserService) CheckEmailAvailability(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) (bool, error) {
	if err := request.Valid(ctx); err != nil {
		return false, errs.NewInvalidValueErr(err.Error())
	}

	emailHash := s.crypto.HashBasic(ctx, []byte(request.Email))

	exists, err := s.repo.ExistsByEmailHash(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// User not found means email is available - this is success
			return true, nil
		case errors.Is(err, errs.ErrConnectionFailure):
			// Database connection issues
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			// Connection pool exhausted
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			// Query was cancelled
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			// Transaction/serialization failure
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			// Database deadlock
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			// Database resources exhausted
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			// Database permission issues
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			// General database error
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			// Malformed query or invalid data
			return false, fmt.Errorf("check email availability: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return false, fmt.Errorf("check email availability cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return false, fmt.Errorf("check email availability timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return false, fmt.Errorf("check email availability: %w", err)
		}
	}

	// Email is available if user does NOT exist
	return !exists, nil
}

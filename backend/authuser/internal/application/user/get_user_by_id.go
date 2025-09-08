package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	// Get user from repository
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// User not found
			return nil, errs.NewNotFoundErr(err, "user by ID")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			// Database connection issues
			return nil, errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			// Query was cancelled
			return nil, fmt.Errorf("get user by ID cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			// Transaction/serialization failure
			return nil, errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			// Database resources exhausted
			return nil, errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			// Database permission issues
			return nil, errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			// General database error
			return nil, errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, errs.ErrInvalidInput):
			// Invalid user ID format
			return nil, errs.NewInvalidValueErr("invalid user ID format")
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return nil, fmt.Errorf("get user by ID cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return nil, fmt.Errorf("get user by ID timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get user by ID: %w", err))
		}
	}

	// Decrypt user data
	if err := s.crypto.DecryptStruct(ctx, user); err != nil {
		return nil, errs.NewNotDecryptedErr("user by ID", err)
	}

	// Convert to response format
	response := user.ToResponse()
	return response, nil
}

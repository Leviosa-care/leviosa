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
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrConnectionFailure):
			// Database connection issues
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			// Connection pool exhausted
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			// Query was cancelled
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			// Transaction/serialization failure
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			// Database deadlock
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			// Database resources exhausted
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			// Database permission issues
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			// General database error
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			// Invalid user ID format
			return nil, fmt.Errorf("get user by ID: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return nil, fmt.Errorf("get user by ID cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return nil, fmt.Errorf("get user by ID timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return nil, fmt.Errorf("get user by ID: %w", err)
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

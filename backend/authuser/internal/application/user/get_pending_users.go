package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (s *UserService) GetPendingUsers(ctx context.Context) ([]*domain.UserResponse, error) {
	// Get pending users from repository
	users, err := s.repo.GetPendingUsers(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// No pending users found - return empty list
			return []*domain.UserResponse{}, nil
		case errors.Is(err, errs.ErrConnectionFailure):
			// Database connection issues
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			// Connection pool exhausted
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			// Query was cancelled
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			// Transaction/serialization failure
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrDeadlock):
			// Database deadlock
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			// Database resources exhausted
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrPermissionDenied):
			// Database permission issues
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			// General database error
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			// Malformed query or invalid data
			return nil, fmt.Errorf("get pending users: %w", err)
		case errors.Is(err, context.Canceled):
			// Request was cancelled
			return nil, fmt.Errorf("get pending users cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			return nil, fmt.Errorf("get pending users timeout: %w", err)
		default:
			// Any unhandled error - wrap with operation context
			return nil, fmt.Errorf("get pending users: %w", err)
		}
	}

	// Decrypt and convert each user to UserResponse
	responses := make([]*domain.UserResponse, 0, len(users))
	for _, user := range users {
		// Decrypt user data
		if err := s.crypto.DecryptStruct(ctx, user); err != nil {
			return nil, errs.NewNotDecryptedErr("pending users list", err)
		}

		// Convert to response format
		response := user.ToResponse()
		responses = append(responses, response)
	}

	return responses, nil
}

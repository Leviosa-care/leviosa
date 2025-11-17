package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *UserService) GetPendingUsers(ctx context.Context) ([]*domain.UserResponse, error) {
	// Get pending users from repository
	usersEncx, err := s.repo.GetPendingUsers(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// No pending users found - return empty list
			return []*domain.UserResponse{}, nil
		}
		return nil, fmt.Errorf("get pending users: %w", err)
	}

	// Decrypt and convert each user to UserResponse
	responses := make([]*domain.UserResponse, 0, len(usersEncx))
	for _, userEncx := range usersEncx {
		// Decrypt user data using the new generated function
		user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("pending users list", err)
		}

		// Convert to response format
		response := user.ToResponse()
		responses = append(responses, response)
	}

	return responses, nil
}

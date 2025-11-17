package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.UserResponse, error) {
	// Get all users from repository
	usersEncx, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		// Only branch on ErrRepositoryNotFound - it has different business logic (empty list is success)
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// No users found - return empty list (success case)
			return []*domain.UserResponse{}, nil
		}
		return nil, fmt.Errorf("get all users: %w", err)
	}

	// Decrypt and convert each user to UserResponse
	responses := make([]*domain.UserResponse, 0, len(usersEncx))
	for _, userEncx := range usersEncx {
		// Decrypt user data using the new generated function
		user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("all users list", err)
		}

		// Convert to response format
		response := user.ToResponse()
		responses = append(responses, response)
	}

	return responses, nil
}

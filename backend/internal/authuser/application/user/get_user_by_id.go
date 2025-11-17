package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	// Get user from repository
	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	// Decrypt user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user by ID", err)
	}

	// Convert to response format
	response := user.ToResponse()
	return response, nil
}

package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *UserService) GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserResponse, error) {
	if appleID == "" {
		return nil, errs.NewInvalidValueErr("Apple ID is required")
	}

	// We need to pass the encrypted Apple ID to match what's stored in DB
	// The repository method expects the encrypted value to match against apple_id_encrypted column
	userEncx, err := s.repo.GetUserByAppleID(ctx, appleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Apple ID: %w", err)
	}

	// Decrypt user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user retrieved by Apple ID", err)
	}

	return user.ToResponse(), nil
}


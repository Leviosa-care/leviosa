package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *UserService) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Validate the new password
	if err := domain.ValidatePassword(newPassword); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by ID for password reset: %w", err)
	}

	// Decrypt user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("user for password reset", err)
	}

	// Update only the password field (no old password verification needed for reset)
	user.Password = newPassword

	// Encrypt the user data using the new generated function
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return errs.NewNotEncryptedErr("user for password reset", err)
	}

	if err := s.repo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return fmt.Errorf("update user password reset: %w", err)
	}

	return nil
}

package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, request *domain.ChangePasswordRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by ID for password change: %w", err)
	}

	// Decrypt the user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("user for password change", err)
	}

	// Verify the old password
	if err := s.VerifyUserPassword(ctx, userID, request.OldPassword); err != nil {
		return fmt.Errorf("verify old password: %w", err)
	}

	// Update only the password field
	user.Password = request.NewPassword

	// Encrypt the user data using the new generated function
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return errs.NewNotEncryptedErr("user for password change", err)
	}

	if err := s.repo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return fmt.Errorf("update user password: %w", err)
	}

	return nil
}

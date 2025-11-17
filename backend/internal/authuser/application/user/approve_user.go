package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *UserService) ApproveUser(ctx context.Context, request *domain.ApproveUserRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	userEncx, err := s.repo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user for approval: %w", err)
	}

	// Decrypt the user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("approved user", err)
	}

	// Verify user is in pending state before approval
	if user.State != domain.Pending {
		return errs.NewConflictErr(fmt.Errorf("user is not in pending state: %s", user.State))
	}

	user.Role = request.Role
	user.State = domain.Active

	// Encrypt the user data using the new generated function
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return errs.NewNotEncryptedErr("approved user", err)
	}

	if err := s.repo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return fmt.Errorf("failed to update user for approval: %w", err)
	}

	return nil
}

package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *UserService) UpdateUserRole(ctx context.Context, request *domain.UpdateUserRoleRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	userEncx, err := s.repo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return fmt.Errorf("get user for role update: %w", err)
	}

	// Decrypt user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("user for role update", err)
	}

	// Update only the role field
	user.Role = request.Role

	// Encrypt the user data using the new generated function
	updatedUserEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return errs.NewNotEncryptedErr("user for role update", err)
	}

	if err := s.repo.UpdateUser(ctx, updatedUserEncx); err != nil {
		return fmt.Errorf("update user role: %w", err)
	}

	return nil
}

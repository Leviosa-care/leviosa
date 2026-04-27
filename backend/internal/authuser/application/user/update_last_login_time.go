package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *UserService) UpdateLastLoginTime(ctx context.Context, userID uuid.UUID) error {
	originalEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by ID to update the LoggedInAt field: %w", err)
	}

	user, err := domain.DecryptUserEncx(ctx, s.crypto, originalEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("user for updated LoggedInAt field: %w", err)

	}

	user.LoggedInAt = time.Now()

	userEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return errs.NewNotEncryptedErr("user for updated LoggedInAt field: %w", err)
	}

	// DecryptUserEncx cannot restore Password from its hash, so ProcessUserEncx
	// would re-hash an empty string. Preserve the original hash instead.
	userEncx.PasswordHashSecure = originalEncx.PasswordHashSecure

	if err := s.repo.UpdateUser(ctx, userEncx); err != nil {
		return fmt.Errorf("update last login time for user: %w", err)
	}

	return nil
}

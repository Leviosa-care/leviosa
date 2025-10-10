package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *UserService) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Validate the new password
	if err := domain.ValidatePassword(newPassword); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("get user for password reset cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("get user for password reset: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("get user for password reset: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("get user for password reset cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("get user for password reset timeout: %w", err)
		default:
			return fmt.Errorf("failed to get user for password reset: %w", err)
		}
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return errs.NewNotUpdatedErr(err, "user")
		case errors.Is(err, errs.ErrUniqueViolation):
			return errs.NewConflictErr(fmt.Errorf("user password reset conflict: %w", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("foreign key constraint violation during password reset: %s", err.Error()))
		case errors.Is(err, errs.ErrNotNullViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("required field missing during password reset: %s", err.Error()))
		case errors.Is(err, errs.ErrCheckViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("data validation failed during password reset: %s", err.Error()))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("update user password reset cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("update user password reset: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("update user password reset: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user data for password reset: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("update user password reset cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("update user password reset timeout: %w", err)
		default:
			return fmt.Errorf("failed to update user password reset: %w", err)
		}
	}

	return nil
}


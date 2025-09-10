package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, request *domain.ChangePasswordRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("get user for password change cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("get user for password change: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("get user for password change: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("get user for password change cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("get user for password change timeout: %w", err)
		default:
			return fmt.Errorf("failed to get user for password change: %w", err)
		}
	}

	if err := s.crypto.DecryptStruct(ctx, user); err != nil {
		return errs.NewNotDecryptedErr("user for password change", err)
	}

	// Verify the old password
	if err := s.VerifyUserPassword(ctx, userID, request.OldPassword); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return errs.NewInvalidValueErr("old password verification failed")
		case errors.Is(err, errs.ErrDomainNotFound):
			return errs.NewNotFoundErr(err, "user")
		default:
			return fmt.Errorf("failed to verify old password: %w", err)
		}
	}

	// Update only the password field
	user.Password = request.NewPassword
	user.PasswordHash = ""

	if err := s.crypto.ProcessStruct(ctx, user); err != nil {
		return errs.NewNotEncryptedErr("user for password change", err)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return errs.NewNotUpdatedErr(err, "user")
		case errors.Is(err, errs.ErrUniqueViolation):
			return errs.NewConflictErr(fmt.Errorf("user password change conflict: %w", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("foreign key constraint violation during password change: %s", err.Error()))
		case errors.Is(err, errs.ErrNotNullViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("required field missing during password change: %s", err.Error()))
		case errors.Is(err, errs.ErrCheckViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("data validation failed during password change: %s", err.Error()))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("update user password cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("update user password: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("update user password: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user data for password change: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("update user password cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("update user password timeout: %w", err)
		default:
			return fmt.Errorf("failed to update user password: %w", err)
		}
	}

	return nil
}


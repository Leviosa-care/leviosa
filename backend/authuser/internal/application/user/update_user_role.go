package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (s *UserService) UpdateUserRole(ctx context.Context, request *domain.UpdateUserRoleRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	user, err := s.repo.GetUserByID(ctx, request.UserID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("get user for role update cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("get user for role update: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("get user for role update: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("get user for role update cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("get user for role update timeout: %w", err)
		default:
			return fmt.Errorf("failed to get user for role update: %w", err)
		}
	}

	if err := s.crypto.DecryptStruct(ctx, user); err != nil {
		return errs.NewNotDecryptedErr("user for role update", err)
	}

	// Update only the role field
	user.Role = request.Role
	user.RoleEncrypted = nil // Clear encrypted field so it gets re-encrypted

	if err := s.crypto.ProcessStruct(ctx, user); err != nil {
		return errs.NewNotEncryptedErr("user for role update", err)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return errs.NewNotUpdatedErr(err, "user")
		case errors.Is(err, errs.ErrUniqueViolation):
			return errs.NewConflictErr(fmt.Errorf("user role update conflict: %w", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("foreign key constraint violation during role update: %s", err.Error()))
		case errors.Is(err, errs.ErrNotNullViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("required field missing during role update: %s", err.Error()))
		case errors.Is(err, errs.ErrCheckViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("data validation failed during role update: %s", err.Error()))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("update user role cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("update user role: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("update user role: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user data for role update: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("update user role cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("update user role timeout: %w", err)
		default:
			return fmt.Errorf("failed to update user role: %w", err)
		}
	}

	return nil
}


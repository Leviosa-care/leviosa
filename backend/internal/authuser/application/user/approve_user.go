package user

import (
	"context"
	"errors"
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("get user for approval cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("get user for approval: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("get user for approval: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user ID: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("get user for approval cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("get user for approval timeout: %w", err)
		default:
			return fmt.Errorf("failed to get user for approval: %w", err)
		}
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return errs.NewNotUpdatedErr(err, "user")
		case errors.Is(err, errs.ErrUniqueViolation):
			return errs.NewConflictErr(fmt.Errorf("user approval conflict: %w", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("foreign key constraint violation during user approval: %s", err.Error()))
		case errors.Is(err, errs.ErrNotNullViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("required field missing during user approval: %s", err.Error()))
		case errors.Is(err, errs.ErrCheckViolation):
			return errs.NewInvalidValueErr(fmt.Sprintf("data validation failed during user approval: %s", err.Error()))
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("update user for approval cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("update user for approval: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("update user for approval: %w", err)
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid user data for approval: %s", err.Error()))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("update user for approval cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("update user for approval timeout: %w", err)
		default:
			return fmt.Errorf("failed to update user for approval: %w", err)
		}
	}

	return nil
}

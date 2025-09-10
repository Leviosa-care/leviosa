package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, request *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	// Get existing user from repository
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get user for update cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return nil, errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.NewInvalidValueErr("invalid user ID format")
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("get user for update cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("get user for update timeout: %w", err)
		default:
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get user for update: %w", err))
		}
	}

	// Decrypt user data to allow field updates
	if err := s.crypto.DecryptStruct(ctx, user); err != nil {
		return nil, errs.NewNotDecryptedErr("user for update", err)
	}

	// Update only non-nil fields from request
	if request.Picture != nil {
		user.Picture = *request.Picture
	}
	if request.FirstName != nil {
		user.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		user.LastName = *request.LastName
	}
	if request.BirthDate != nil {
		user.BirthDate = *request.BirthDate
	}
	if request.Gender != nil {
		user.Gender = *request.Gender
	}
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Telephone != nil {
		user.Telephone = *request.Telephone
	}
	if request.PostalCode != nil {
		user.PostalCode = *request.PostalCode
	}
	if request.City != nil {
		user.City = *request.City
	}
	if request.Address1 != nil {
		user.Address1 = *request.Address1
	}
	if request.Address2 != nil {
		user.Address2 = *request.Address2
	}

	// Encrypt user data before saving
	if err := s.crypto.ProcessStruct(ctx, user); err != nil {
		return nil, errs.NewNotEncryptedErr("user for update", err)
	}

	// Save updated user to repository
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return nil, errs.NewNotUpdatedErr(err, "user")
		case errors.Is(err, errs.ErrUniqueViolation):
			return nil, errs.NewAlreadyExistsError(err, "user with this email or phone")
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return nil, errs.NewInvalidValueErr("invalid reference in user data")
		case errors.Is(err, errs.ErrNotNullViolation):
			return nil, errs.NewInvalidValueErr("required field is missing")
		case errors.Is(err, errs.ErrCheckViolation):
			return nil, errs.NewInvalidValueErr("user data violates database constraints")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("update user cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return nil, errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("update user cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("update user timeout: %w", err)
		default:
			return nil, errs.NewInternalErr(fmt.Errorf("failed to update user: %w", err))
		}
	}

	// Decrypt user data again for response
	if err := s.crypto.DecryptStruct(ctx, user); err != nil {
		return nil, errs.NewNotDecryptedErr("updated user", err)
	}

	// Convert to response format
	response := user.ToResponse()
	return response, nil
}


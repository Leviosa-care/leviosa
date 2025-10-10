package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/validation"
	"github.com/google/uuid"
)

func (s *UserService) CreatePendingUser(ctx context.Context, email string) (uuid.UUID, error) {
	if err := validation.ValidateEmail(email); err != nil {
		return uuid.Nil, errs.NewInvalidValueErr(err.Error())
	}

	emailHash := s.crypto.HashBasic(ctx, []byte(email))

	// Check if user already exists
	existingUserEncx, err := s.repo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// User doesn't exist, proceed with creation
			user := newPendingUser(email)

			// Encrypt the user data using the new generated function
			userEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
			if err != nil {
				return uuid.Nil, errs.NewInternalErr(fmt.Errorf("failed to encrypt user data: %w", err))
			}

			// Create user in database
			if err := s.repo.CreateUser(ctx, userEncx); err != nil {
				switch {
				case errors.Is(err, errs.ErrUniqueViolation):
					return uuid.Nil, errs.NewConflictErr(errors.New("user already exists"))
				case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
					return uuid.Nil, errs.NewExternalServiceErr(err, "database unavailable")
				case errors.Is(err, errs.ErrInvalidInput):
					return uuid.Nil, errs.NewInvalidValueErr("invalid user data")
				default:
					return uuid.Nil, errs.NewInternalErr(fmt.Errorf("failed to create user: %w", err))
				}
			}

			return userEncx.ID, nil
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return uuid.Nil, errs.NewExternalServiceErr(err, "database unavailable")
		default:
			return uuid.Nil, errs.NewInternalErr(fmt.Errorf("failed to check user existence: %w", err))
		}
	}

	// User already exists, return the existing user's ID with conflict error
	// This allows callers to still get the user ID even when there's a conflict
	return existingUserEncx.ID, errs.NewConflictErr(errors.New("user already exists"))
}

func newPendingUser(email string) *domain.User {
	return &domain.User{
		ID:    uuid.New(),
		Email: email,
		State: domain.Unverified,
	}
}

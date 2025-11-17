package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func (s *UserService) CreatePendingUser(ctx context.Context, email string) (uuid.UUID, error) {
	if err := validation.ValidateEmail(email); err != nil {
		return uuid.Nil, errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return uuid.Nil, errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	// Check if user already exists
	existingUserEncx, err := s.repo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		// Only branch on ErrRepositoryNotFound - it triggers different business logic
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// User doesn't exist, CREATE ONE
			user := newPendingUser(email)

			userEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
			if err != nil {
				return uuid.Nil, fmt.Errorf("encrypt user data: %w", err)
			}

			// Create user in database
			if err := s.repo.CreateUser(ctx, userEncx); err != nil {
				return uuid.Nil, fmt.Errorf("create user: %w", err)
			}

			return userEncx.ID, nil
		}

		// All other errors - just wrap and return
		return uuid.Nil, fmt.Errorf("check user existence: %w", err)
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

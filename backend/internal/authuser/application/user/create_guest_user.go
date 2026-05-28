package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// CreateGuestUser creates a fully active user account from the minimal data
// captured during a guest booking (name, email, phone, password).
// The account has profile_incomplete = true, meaning gender, birthdate, and
// address are empty and should be prompted later via an in-app nudge.
func (s *UserService) CreateGuestUser(ctx context.Context, req *domain.GuestClaimRequest) (uuid.UUID, error) {
	if err := req.Valid(ctx); err != nil {
		return uuid.Nil, errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(req.Email)
	if err != nil {
		return uuid.Nil, errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	// Check if user already exists
	_, err = s.repo.GetUserByEmailHash(ctx, emailHash)
	if err == nil {
		// User already exists
		return uuid.Nil, errs.NewConflictErr(errors.New("email is already registered"))
	}
	if !errors.Is(err, errs.ErrRepositoryNotFound) {
		return uuid.Nil, fmt.Errorf("check user existence: %w", err)
	}

	// Create the guest user with minimal profile
	user := &domain.User{
		ID:               uuid.New(),
		Email:            req.Email,
		Password:         req.Password,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Telephone:        req.Phone,
		State:            domain.Active,
		Role:             "standard",
		ProfileIncomplete: true,
		CreatedAt:        time.Now(),
		LoggedInAt:       time.Now(),
	}

	userEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("encrypt user data: %w", err)
	}

	if err := s.repo.CreateUser(ctx, userEncx); err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return userEncx.ID, nil
}

// ValidateGuestClaimEmail checks if the email is available for guest claim.
// Returns nil if available, conflict error if already registered.
func (s *UserService) ValidateGuestClaimEmail(ctx context.Context, email string) error {
	if err := validation.ValidateEmail(email); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	exists, err := s.repo.ExistsByEmailHash(ctx, emailHash)
	if err != nil {
		return fmt.Errorf("check email existence: %w", err)
	}
	if exists {
		return errs.NewConflictErr(errors.New("email is already registered"))
	}

	return nil
}

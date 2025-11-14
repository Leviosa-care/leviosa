package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *UserService) VerifyUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	// Get the user from repository to access the stored password hash
	userEncx, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to get user: %w", err))
		}
	}

	ok, err := s.crypto.CompareSecureHashAndValue(ctx, password, userEncx.PasswordHashSecure)
	if err != nil {
		return errs.NewUnexpectedError(err)
	}
	if !ok {
		return errs.NewInvalidValueErr("password verification failed: provided password does not match stored hash")
	}

	return nil
}

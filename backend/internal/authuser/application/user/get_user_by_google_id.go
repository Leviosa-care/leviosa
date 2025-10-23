package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *UserService) GetUserByGoogleID(ctx context.Context, googleID string) (*domain.UserResponse, error) {
	if googleID == "" {
		return nil, errs.NewInvalidValueErr("Google ID is required")
	}

	// We need to pass the encrypted Google ID to match what's stored in DB
	// The repository method expects the encrypted value to match against google_id_encrypted column
	userEncx, err := s.repo.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "user")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, errs.NewExternalServiceErr(err, "database query cancelled")
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return nil, errs.NewExternalServiceErr(err, "database transaction failed")
		default:
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get user by Google ID: %w", err))
		}
	}

	// Decrypt user data using the new generated function
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user retrieved by Google ID", err)
	}

	return user.ToResponse(), nil
}


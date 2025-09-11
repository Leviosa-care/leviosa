package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (s *UserService) GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserResponse, error) {
	if appleID == "" {
		return nil, errs.NewInvalidValueErr("Apple ID is required")
	}

	// We need to pass the encrypted Apple ID to match what's stored in DB
	// The repository method expects the encrypted value to match against apple_id_encrypted column
	user, err := s.repo.GetUserByAppleID(ctx, appleID)
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
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get user by Apple ID: %w", err))
		}
	}

	if err := s.crypto.DecryptStruct(ctx, user); err != nil {
		return nil, errs.NewNotDecryptedErr("user retrieved by Apple ID", err)
	}

	return user.ToResponse(), nil
}
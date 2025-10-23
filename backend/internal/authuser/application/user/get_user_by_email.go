package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/hengadev/encx"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	if err := validation.ValidateEmail(email); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	userEncx, err := s.repo.GetUserByEmailHash(ctx, emailHash)
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
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get user by email: %w", err))
		}
	}

	// Decrypt user data
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user retrieved by email", err)
	}

	return user.ToResponse(), nil
}

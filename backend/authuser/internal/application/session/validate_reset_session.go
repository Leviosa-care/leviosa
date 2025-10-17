package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/core/errs"
	"github.com/hengadev/encx"
)

func (s *SessionService) ValidateResetSession(ctx context.Context, token string) (string, error) {
	// Hash the token for lookup
	tokenBytes, err := encx.SerializeValue(token)
	if err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}
	tokenHash := s.crypto.HashBasic(ctx, tokenBytes)

	userEmailHash, err := s.repo.ValidateResetSession(ctx, tokenHash)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return "", errs.NewNotFoundErr(err, "reset session")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return "", errs.NewExternalServiceErr(err, "Redis unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return "", errs.NewExternalServiceErr(err, "Redis resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return "", fmt.Errorf("validate reset session cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return "", errs.NewExternalServiceErr(err, "Redis transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return "", fmt.Errorf("validate reset session: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return "", fmt.Errorf("validate reset session: %w", err)
		case errors.Is(err, context.Canceled):
			return "", fmt.Errorf("validate reset session cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return "", fmt.Errorf("validate reset session timeout: %w", err)
		default:
			return "", fmt.Errorf("failed to validate reset session: %w", err)
		}
	}

	return userEmailHash, nil
}

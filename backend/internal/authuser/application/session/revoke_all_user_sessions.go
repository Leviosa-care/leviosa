package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func (s *SessionService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	userIDBytes, err := encx.SerializeValue(userID)
	if err != nil {
		return errs.NewInvalidValueErr(fmt.Sprintf("failed to serialize userID: %w", err))
	}
	userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

	// Call repository with the hashed userID
	if err := s.repo.RevokeAllUserSessions(ctx, userIDHash); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// No sessions found for user - this is actually success for revocation
			return nil
		case errors.Is(err, errs.ErrConnectionFailure):
			return errs.NewExternalServiceErr(err, "session service unavailable during revocation")
		case errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "session service overloaded during revocation")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "session service resources exhausted during revocation")
		case errors.Is(err, errs.ErrQueryCancelled):
			return errs.NewQueryFailedErr(fmt.Errorf("session revocation query cancelled: %w", err))
		case errors.Is(err, errs.ErrTransactionFailure):
			return errs.NewExternalServiceErr(err, "session revocation transaction failed, retry possible")
		case errors.Is(err, errs.ErrPermissionDenied):
			return errs.NewForbiddenErr("session revocation permission denied")
		case errors.Is(err, errs.ErrContext):
			return errs.NewQueryFailedErr(fmt.Errorf("session revocation request timed out: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database error during session revocation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unexpected error during session revocation for user %s: %w", userID, err))
		}
	}

	return nil
}

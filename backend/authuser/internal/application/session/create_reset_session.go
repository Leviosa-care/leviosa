package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/core/errs"
)

func (s *SessionService) CreateResetSession(ctx context.Context, token, userEmail string, ttl time.Duration) error {
	// Hash the token for storage
	tokenHash := s.crypto.HashBasic(ctx, []byte(token))

	// Store plaintext email directly (no hashing needed)
	if err := s.repo.StoreResetSession(ctx, tokenHash, userEmail, ttl); err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "Redis unavailable")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "Redis resources exhausted")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("store reset session cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return errs.NewExternalServiceErr(err, "Redis transaction failed")
		case errors.Is(err, errs.ErrPermissionDenied):
			return fmt.Errorf("store reset session: %w", err)
		case errors.Is(err, errs.ErrDatabase):
			return fmt.Errorf("store reset session: %w", err)
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("store reset session cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("store reset session timeout: %w", err)
		default:
			return fmt.Errorf("failed to store reset session: %w", err)
		}
	}

	return nil
}


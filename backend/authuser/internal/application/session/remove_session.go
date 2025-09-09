package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
)

func (s *SessionService) RemoveSession(ctx context.Context, sessionID uuid.UUID) error {
	// Find the session using the access token hash
	sessionBytes, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "session to remove")
		case errors.Is(err, errs.ErrConnectionFailure):
			return fmt.Errorf("remove session - database connection failed: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return fmt.Errorf("remove session - database overloaded: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("remove session - operation timed out: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return fmt.Errorf("remove session - database resources exhausted: %w", err)
		default:
			return fmt.Errorf("remove session - failed to find session: %w", err)
		}
	}

	// Decode and decrypt the session
	var session session.Session
	if err := json.Unmarshal(sessionBytes, &session); err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	if err := s.crypto.DecryptStruct(ctx, &session); err != nil {
		return errs.NewNotDecryptedErr("session", err)
	}

	// Remove the session by ID (this removes the session data and both token mappings)
	if err = s.repo.RemoveSessionByID(ctx, session.ID); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "session to remove")
		case errors.Is(err, errs.ErrConnectionFailure):
			return fmt.Errorf("remove session - database connection failed during deletion: %w", err)
		case errors.Is(err, errs.ErrTooManyConnections):
			return fmt.Errorf("remove session - database overloaded during deletion: %w", err)
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("remove session - deletion operation timed out: %w", err)
		case errors.Is(err, errs.ErrResourceExhausted):
			return fmt.Errorf("remove session - database resources exhausted during deletion: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure):
			return fmt.Errorf("remove session - transaction failed during deletion: %w", err)
		default:
			return fmt.Errorf("remove session - failed to remove session: %w", err)
		}
	}

	return nil
}

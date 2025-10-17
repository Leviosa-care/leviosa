package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *SessionService) UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, completedAt *time.Time) error {
	// First get the current session to update it
	sessionData, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(fmt.Errorf("session not found: %w", err), "session")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewExternalServiceErr(err, "database connection error")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	var sessionEncx *session.SessionEncx
	err = json.Unmarshal(sessionData, &sessionEncx)
	if err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	// Decrypt session using the new generated function
	sess, err := session.DecryptSessionEncx(ctx, s.crypto, sessionEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("session retrieved during user completion process", err)
	}

	sess.CompletedAt = completedAt

	// Encrypt updated session using the new generated function
	// BUG: This is the part that is causing the issue since we have *time.Time
	updatedSessionEncx, err := session.ProcessSessionEncx(ctx, s.crypto, sess)
	if err != nil {
		return errs.NewNotEncryptedErr("session during completion update", err)
	}

	updatedSessionData, err := json.Marshal(updatedSessionEncx)
	if err != nil {
		return errs.NewJSONMarshalErr(err)
	}

	// Update the session in the repository
	if err := s.repo.UpdateSessionCompletion(ctx, sessionID, updatedSessionData); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(fmt.Errorf("session not found during update: %w", err), "session")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed during update: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewExternalServiceErr(err, "database connection error during update")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during update: %w", err))
		}
	}

	return nil
}

package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/google/uuid"
)

func (s *SessionService) UpdateSessionCompletion(ctx context.Context, sessionID string, completedAt *time.Time) error {
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid session ID format")
	}

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

	// NOTE: what was that ?
	// // Parse the session data (it's encrypted JSON)
	// var sessionMap map[string]any
	// err = json.Unmarshal(sessionData, &sessionMap)
	// if err != nil {
	// 	return errs.NewJSONUnmarshalErr(err)
	// }
	//
	// // If we have a completed timestamp, encrypt it and add to the session
	// if completedAt != nil {
	// 	// Here we'd normally use the crypto service to encrypt the timestamp
	// 	// For now, we'll add the encrypted fields to the map
	// 	sessionMap["completed_at"] = completedAt
	// 	// In a real implementation, you'd encrypt this and store in completed_at_encrypted
	// }
	// // Re-encode the session
	// updatedSessionData, err := json.Marshal(sessionMap)
	// if err != nil {
	// 	return errs.NewJSONMarshalErr(err)
	// }

	var session *auth.Session
	err = json.Unmarshal(sessionData, &session)
	if err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	if err := s.crypto.DecryptStruct(ctx, session); err != nil {
		return errs.NewNotDecryptedErr("session retrieved during user completion process", err)
	}

	session.CompletedAt = completedAt

	updatedSessionData, err := json.Marshal(session)
	if err != nil {
		return errs.NewJSONMarshalErr(err)
	}

	// Update the session in the repository
	if err := s.repo.UpdateSessionCompletion(ctx, sessionUUID, updatedSessionData); err != nil {
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

package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *SessionService) UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, completedAt *time.Time) error {
	// First get the current session to update it
	sessionData, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("find session by ID for completion update: %w", err)
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
		return fmt.Errorf("update session completion: %w", err)
	}

	return nil
}

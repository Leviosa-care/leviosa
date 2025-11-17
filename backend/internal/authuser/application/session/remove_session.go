package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *SessionService) RemoveSession(ctx context.Context, sessionID uuid.UUID) error {
	// Find the session using the access token hash
	sessionBytes, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("find session by ID for removal: %w", err)
	}

	// Decode and decrypt the session using the new generated function
	var sessionEncx session.SessionEncx
	if err := json.Unmarshal(sessionBytes, &sessionEncx); err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	session, err := session.DecryptSessionEncx(ctx, s.crypto, &sessionEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("session", err)
	}

	// Remove the session by ID (this removes the session data and both token mappings)
	if err = s.repo.RemoveSessionByID(ctx, session.ID); err != nil {
		return fmt.Errorf("remove session by ID: %w", err)
	}

	return nil
}

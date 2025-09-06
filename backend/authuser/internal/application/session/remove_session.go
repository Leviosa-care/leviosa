package session

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
)

func (s *SessionService) RemoveSession(ctx context.Context, sessionID uuid.UUID) error {
	// Find the session using the access token hash
	sessionBytes, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewNotFoundErr(err, "session to remove")
		}
		return err
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
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewNotFoundErr(err, "session to remove")
		}
		return err
	}

	return nil
}

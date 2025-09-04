package session

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (s *SessionService) RemoveSession(ctx context.Context, request *domain.RemoveSessionRequest) error {
	// Hash the access token to find the session
	tokenHash := s.crypto.HashBasic(ctx, []byte(request.Token))

	// Find the session using the access token hash
	sessionID, sessionBytes, err := s.repo.FindSessionByAccessToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewNotFoundErr(err, "session to remove")
		}
		return err
	}

	// Decode and decrypt the session
	var session auth.Session
	if err := json.Unmarshal(sessionBytes, &session); err != nil {
		return errs.NewJSONUnmarshalErr(err)
	}

	if err := s.crypto.DecryptStruct(ctx, &session); err != nil {
		return errs.NewNotDecryptedErr("session", err)
	}

	session.ID, err = uuid.Parse(sessionID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid session ID format")
	}

	// Remove the session by ID (this removes the session data and both token mappings)
	err = s.repo.RemoveSessionByID(ctx, session.ID.String())
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewNotFoundErr(err, "session to remove")
		}
		return err
	}

	return nil
}

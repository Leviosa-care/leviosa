package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (s *SessionService) GetSession(ctx context.Context, request *domain.GetSessionRequest) (*auth.Session, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	tokenHash := s.crypto.HashBasic(ctx, []byte(request.Token))

	sessionID, sessionBytes, err := s.repo.FindSessionByAccessTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewNotFoundErr(err, "session")
		}
		return nil, errs.NewUnexpectedError(err)
	}

	var session auth.Session
	if err := json.Unmarshal(sessionBytes, &session); err != nil {
		return nil, errs.NewJSONUnmarshalErr(err)
	}

	if err := s.crypto.DecryptStruct(ctx, &session); err != nil {
		return nil, errs.NewNotDecryptedErr("session", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errs.NewExpiredTokenErr("session", nil)
	}

	session.ID, err = uuid.Parse(sessionID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid session ID format")
	}

	return &session, nil
}

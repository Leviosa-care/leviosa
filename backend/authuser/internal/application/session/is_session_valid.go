package session

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"time"
//
// 	"github.com/Leviosa-care/authuser/internal/domain"
// 	"github.com/google/uuid"
//
// 	"github.com/Leviosa-care/core/auth/session"
// 	"github.com/Leviosa-care/core/errs"
// )
//
// func (s *SessionService) IsSessionValid(ctx context.Context, request *domain.ValidateSessionRequest) (bool, error) {
// 	if err := request.Valid(ctx); err != nil {
// 		return false, errs.NewInvalidValueErr(err.Error())
// 	}
//
// 	tokenHash := s.crypto.HashBasic(ctx, []byte(request.Token))
//
// 	sessionID, encryptedSession, err := s.repo.FindSessionByAccessTokenHash(ctx, tokenHash)
// 	if err != nil {
// 		if errors.Is(err, errs.ErrRepositoryNotFound) || errors.Is(err, errs.ErrExpiredToken) {
// 			return false, errs.NewNotFoundErr(err, "session by access token for validation")
// 		}
// 		return false, errs.NewUnexpectedError(err)
// 	}
//
// 	var sess session.Session
// 	if err := json.Unmarshal(encryptedSession, &sess); err != nil {
// 		return false, errs.NewJSONUnmarshalErr(err)
// 	}
//
// 	sess.ID, err = uuid.Parse(sessionID)
// 	if err != nil {
// 		return false, errs.NewInvalidValueErr("invalid s ID format")
// 	}
//
// 	return sess.Valid(ctx) == nil && time.Now().Before(sess.ExpiresAt), nil
// }

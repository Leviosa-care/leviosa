package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"

	authsession "github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *SessionService) RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error) {
	// Find session by refresh token
	sessionData, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(fmt.Errorf("session not found during refresh: %w", err), "session")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed during refresh: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during refresh: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during refresh: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during refresh: %w", err))
		}
	}

	// Decode session as SessionEncx
	var sessionEncx authsession.SessionEncx
	if err := json.Unmarshal(sessionData, &sessionEncx); err != nil {
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to decode session during refresh: %w", err))
	}

	// Decrypt session using the new generated function
	session, err := authsession.DecryptSessionEncx(ctx, s.crypto, &sessionEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("refresh session", err)
	}

	session.ID = sessionID

	// Verify session state - both pending and active sessions can be refreshed
	if session.State != authsession.SessionActive && session.State != authsession.SessionPending {
		return nil, errs.NewUnauthorizedErr("invalid session state for refresh")
	}

	// Generate new token pair
	newAccessToken, err := authsession.GenerateToken()
	if err != nil {
		return nil, errs.NewUnexpectedError(err)
	}

	newRefreshToken, err := authsession.GenerateToken()
	if err != nil {
		return nil, errs.NewUnexpectedError(err)
	}

	// Get token durations from cache
	accessDuration := s.cache.GetAccessTokenDuration()
	refreshDuration := s.cache.GetRefreshTokenDuration()

	// Update session with new tokens
	session.AccessToken = newAccessToken
	session.RefreshToken = newRefreshToken

	// Encrypt updated session using the new generated function
	updatedSessionEncx, err := authsession.ProcessSessionEncx(ctx, s.crypto, session)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("updated session", err)
	}

	// Encode updated session data
	updatedSessionData, err := json.Marshal(updatedSessionEncx)
	if err != nil {
		return nil, errs.NewJSONMarshalErr(err)
	}

	now := time.Now()
	accessExpiry := now.Add(accessDuration)
	refreshExpiry := now.Add(refreshDuration)

	// Perform token rotation - replace old tokens with new ones
	if err := s.repo.RefreshTokenPair(
		ctx,
		sessionEncx.RefreshTokenHash,         // oldRefreshTokenHash
		updatedSessionEncx.AccessTokenHash,   // newAccessTokenHash
		updatedSessionEncx.RefreshTokenHash,  // newRefreshTokenHash
		session.ID,                           // sessionID (uuid.UUID)
		updatedSessionData,                   // updatedSessionData
		accessDuration,
		refreshDuration,
	); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(fmt.Errorf("session not found during token rotation: %w", err), "session")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed during token rotation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error during token rotation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during token rotation: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during token rotation: %w", err))
		}
	}

	return &domain.RefreshSessionResponse{
		AccessToken:        newAccessToken,
		RefreshToken:       newRefreshToken,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}, nil
}

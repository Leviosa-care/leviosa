package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	authsession "github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *SessionService) RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error) {
	// Find session by refresh token
	sessionData, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("find session by ID for refresh: %w", err)
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
		sessionEncx.RefreshTokenHash,        // oldRefreshTokenHash
		updatedSessionEncx.AccessTokenHash,  // newAccessTokenHash
		updatedSessionEncx.RefreshTokenHash, // newRefreshTokenHash
		session.ID,                          // sessionID (uuid.UUID)
		updatedSessionData,                  // updatedSessionData
		accessDuration,
		refreshDuration,
	); err != nil {
		return nil, fmt.Errorf("refresh token pair: %w", err)
	}

	return &domain.RefreshSessionResponse{
		AccessToken:        newAccessToken,
		RefreshToken:       newRefreshToken,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}, nil
}

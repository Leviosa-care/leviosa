package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *SessionService) CreateSession(ctx context.Context, request *domain.CreateSessionRequest) (*domain.CreateSessionResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	userID, _ := uuid.Parse(request.UserID)
	now := time.Now()

	// Generate access and refresh tokens
	accessToken, err := session.GenerateToken()
	if err != nil {
		return nil, errs.NewUnexpectedError(err)
	}

	refreshToken, err := session.GenerateToken()
	if err != nil {
		return nil, errs.NewUnexpectedError(err)
	}

	// Get token durations from cache
	accessDuration := s.cache.GetAccessTokenDuration()
	refreshDuration := s.cache.GetRefreshTokenDuration()

	// Use shorter durations for pending sessions
	if request.State == session.SessionPending {
		accessDuration = session.PendingSessionDuration
		refreshDuration = session.PendingSessionDuration
	}

	sess := &session.Session{
		ID:           uuid.New(),
		UserID:       userID,
		Role:         request.Role,
		State:        request.State,
		CreatedAt:    now,
		ExpiresAt:    now.Add(refreshDuration), // Session expires with refresh token
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Encrypt session and token pair using the new generated function
	sessionEncx, err := session.ProcessSessionEncx(ctx, s.crypto, sess)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("create session", err)
	}

	sessionEncoded, err := json.Marshal(sessionEncx)
	if err != nil {
		return nil, errs.NewJSONMarshalErr(err)
	}

	accessExpiry := now.Add(accessDuration)
	refreshExpiry := now.Add(refreshDuration)

	if err := s.repo.CreateSession(ctx, sessionEncx.ID, sessionEncx.AccessTokenHash, sessionEncx.RefreshTokenHash, sessionEncx.UserIDHash, sessionEncoded, accessDuration, refreshDuration); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &domain.CreateSessionResponse{
		RefreshToken:       sess.RefreshToken, // These are still accessible from the original session
		AccessToken:        sess.AccessToken,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}, nil
}

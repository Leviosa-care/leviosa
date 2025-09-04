package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/google/uuid"
)

func (s *SessionService) RefreshSession(ctx context.Context, request *domain.RefreshSessionRequest) (*domain.RefreshSessionResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Find session by refresh token
	sessionID, sessionData, err := s.repo.FindSessionByRefreshToken(ctx, request.RefreshToken)
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

	// Decode session
	session, err := auth.DecodeSession(sessionData)
	if err != nil {
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to decode session during refresh: %w", err))
	}

	// Verify session state - both pending and active sessions can be refreshed
	if session.State != auth.SessionActive && session.State != auth.SessionPending {
		return nil, errs.NewUnauthorizedErr("invalid session state for refresh")
	}

	session.ID, err = uuid.Parse(sessionID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid session ID format")
	}

	// Generate new token pair
	newAccessToken, err := auth.GenerateToken()
	if err != nil {
		return nil, errs.NewUnexpectedError(err)
	}

	newRefreshToken, err := auth.GenerateToken()
	if err != nil {
		return nil, errs.NewUnexpectedError(err)
	}

	// Get token durations from cache
	accessDuration := s.cache.GetAccessTokenDuration()
	refreshDuration := s.cache.GetRefreshTokenDuration()

	// Create new token pair
	newTokenPair := &auth.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	// Encrypt new token pair
	s.crypto.ProcessStruct(ctx, newTokenPair)

	now := time.Now()
	accessExpiry := now.Add(accessDuration)
	refreshExpiry := now.Add(refreshDuration)

	// Perform token rotation - replace old tokens with new ones
	if err := s.repo.RefreshTokenPair(
		ctx,
		request.RefreshToken,          // oldRefreshTokenHash
		newTokenPair.AccessTokenHash,  // newAccessTokenHash
		newTokenPair.RefreshTokenHash, // newRefreshTokenHash
		session.ID,                    // sessionID (uuid.UUID)
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
		AccessToken:        newTokenPair.AccessTokenHash,
		RefreshToken:       newTokenPair.RefreshTokenHash,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}, nil
}

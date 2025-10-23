package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// CreateSession creates session data with both access and refresh tokens
// Implements secure two-step lookup: token -> session ID -> session data
func (r *SessionRepository) CreateSession(ctx context.Context, sessionID uuid.UUID, accessTokenHash, refreshTokenHash, userIDHash string, sessionEncoded []byte, accessTTL, refreshTTL time.Duration) error {
	sessionKey := session.FormatSessionKey(sessionID.String())
	accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(refreshTokenHash)
	userSessionIndexKey := session.FormatUserSessionIndexKey(userIDHash)
	sessionIDStr := sessionID.String()

	var creationErrs errsx.Map

	// Store session data with the longer TTL (refresh token duration)
	if err := r.client.Set(ctx, sessionKey, sessionEncoded, refreshTTL).Err(); err != nil {
		creationErrs.Set("session_data", errs.ClassifyRedisError("save session data", err).Error())
		return creationErrs.AsError()
	}

	// Store access token -> session ID mapping
	if err := r.client.Set(ctx, accessTokenKey, sessionIDStr, accessTTL).Err(); err != nil {
		// Rollback session data
		if delErr := r.client.Del(ctx, sessionKey).Err(); delErr != nil {
			creationErrs.Set("rollback_session", delErr.Error())
		}
		creationErrs.Set("access_token", errs.ClassifyRedisError("save access token mapping", err).Error())
		return creationErrs.AsError()
	}

	// Store refresh token -> session ID mapping
	if err := r.client.Set(ctx, refreshTokenKey, sessionIDStr, refreshTTL).Err(); err != nil {
		// Rollback both session data and access token
		if delErr := r.client.Del(ctx, sessionKey).Err(); delErr != nil {
			creationErrs.Set("rollback_session", delErr.Error())
		}
		if delErr := r.client.Del(ctx, accessTokenKey).Err(); delErr != nil {
			creationErrs.Set("rollback_access_token", delErr.Error())
		}
		creationErrs.Set("refresh_token", errs.ClassifyRedisError("save refresh token mapping", err).Error())
		return creationErrs.AsError()
	}

	// Add session ID to user session index
	if err := r.client.SAdd(ctx, userSessionIndexKey, sessionIDStr).Err(); err != nil {
		// Rollback all previous operations
		if delErr := r.client.Del(ctx, sessionKey).Err(); delErr != nil {
			creationErrs.Set("rollback_session", delErr.Error())
		}
		if delErr := r.client.Del(ctx, accessTokenKey).Err(); delErr != nil {
			creationErrs.Set("rollback_access_token", delErr.Error())
		}
		if delErr := r.client.Del(ctx, refreshTokenKey).Err(); delErr != nil {
			creationErrs.Set("rollback_refresh_token", delErr.Error())
		}
		creationErrs.Set("user_session_index", errs.ClassifyRedisError("add to user session index", err).Error())
		return creationErrs.AsError()
	}

	return nil
}

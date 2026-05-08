package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// RefreshTokenPair rotates both the access and refresh tokens for a session.
// The old tokens are deleted directly by their known hashes (O(1)), avoiding a SCAN.
func (r *SessionRepository) RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, oldAccessTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, updatedSessionData []byte, accessTTL, refreshTTL time.Duration) error {
	sessionIDStr := sessionID.String()
	oldAccessTokenKey := session.FormatAccessTokenKey(oldAccessTokenHash)
	oldRefreshTokenKey := session.FormatRefreshTokenKey(oldRefreshTokenHash)
	newAccessTokenKey := session.FormatAccessTokenKey(newAccessTokenHash)
	newRefreshTokenKey := session.FormatRefreshTokenKey(newRefreshTokenHash)

	var refreshErrs errsx.Map

	// Create new access token
	if err := r.client.Set(ctx, newAccessTokenKey, sessionIDStr, accessTTL).Err(); err != nil {
		refreshErrs.Set("new_access_token", errs.ClassifyRedisError("save new access token", err).Error())
		return refreshErrs.AsError()
	}

	// Create new refresh token
	if err := r.client.Set(ctx, newRefreshTokenKey, sessionIDStr, refreshTTL).Err(); err != nil {
		if delErr := r.client.Del(ctx, newAccessTokenKey).Err(); delErr != nil {
			refreshErrs.Set("rollback_new_access_token", delErr.Error())
		}
		refreshErrs.Set("new_refresh_token", errs.ClassifyRedisError("save new refresh token", err).Error())
		return refreshErrs.AsError()
	}

	// Update session data with new token hashes
	sessionKey := session.FormatSessionKey(sessionIDStr)
	if err := r.client.Set(ctx, sessionKey, updatedSessionData, 0).Err(); err != nil {
		if delErr := r.client.Del(ctx, newAccessTokenKey, newRefreshTokenKey).Err(); delErr != nil {
			refreshErrs.Set("rollback_new_tokens", delErr.Error())
		}
		refreshErrs.Set("update_session_data", errs.ClassifyRedisError("update session data", err).Error())
		return refreshErrs.AsError()
	}

	// Delete old tokens. Failures here are non-fatal — the new tokens are already active
	// and the old ones will expire on their own TTL.
	if err := r.client.Del(ctx, oldAccessTokenKey, oldRefreshTokenKey).Err(); err != nil {
		logger, logErr := ctxutil.GetLoggerFromContext(ctx)
		if logErr != nil {
			return nil
		}
		logger.WarnContext(ctx, "Session: Failed to delete old tokens after rotation",
			"operation", "refresh_token_pair_cleanup",
			"error", err)
	}

	return nil
}

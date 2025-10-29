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

// RefreshTokenPair refreshes access token and rotates refresh token
// This method cleans up ALL existing tokens for the session to ensure security
func (r *SessionRepository) RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, updatedSessionData []byte, accessTTL, refreshTTL time.Duration) error {
	sessionIDStr := sessionID.String()
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
		// Rollback new access token
		if delErr := r.client.Del(ctx, newAccessTokenKey).Err(); delErr != nil {
			refreshErrs.Set("rollback_new_access_token", delErr.Error())
		}
		refreshErrs.Set("new_refresh_token", errs.ClassifyRedisError("save new refresh token", err).Error())
		return refreshErrs.AsError()
	}

	// Update session data with new token hashes
	sessionKey := session.FormatSessionKey(sessionIDStr)
	if err := r.client.Set(ctx, sessionKey, updatedSessionData, 0).Err(); err != nil {
		// Rollback new tokens
		if delErr := r.client.Del(ctx, newAccessTokenKey, newRefreshTokenKey).Err(); delErr != nil {
			refreshErrs.Set("rollback_new_tokens", delErr.Error())
		}
		refreshErrs.Set("update_session_data", errs.ClassifyRedisError("update session data", err).Error())
		return refreshErrs.AsError()
	}

	// Remove old refresh token
	if err := r.client.Del(ctx, oldRefreshTokenKey).Err(); err != nil {
		// Don't fail the operation if old token deletion fails, just log
		refreshErrs.Set("cleanup_old_refresh_token", err.Error())
	}

	// Clean up any stale access tokens for this session by scanning and removing them
	// This ensures security by invalidating all old access tokens
	pattern := session.FormatAccessTokenKey("*")
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if key == newAccessTokenKey {
			// Skip the new access token we just created
			continue
		}

		// Check if this access token maps to our session ID
		storedSessionID, err := r.client.Get(ctx, key).Result()
		if err != nil {
			// Key might have expired or been deleted, skip it
			continue
		}

		if storedSessionID == sessionIDStr {
			// This is an old access token for our session, remove it
			if err := r.client.Del(ctx, key).Err(); err != nil {
				refreshErrs.Set("cleanup_stale_access_token", err.Error())
			}
		}
	}

	if err := iter.Err(); err != nil {
		refreshErrs.Set("scan_access_tokens", err.Error())
	}

	// If there were any cleanup errors, log them but don't fail the operation
	if !refreshErrs.IsEmpty() {
		logger, err := ctxutil.GetLoggerFromContext(ctx)
		if err != nil {
			return err
		}
		logger.WarnContext(ctx, "Session: Token pair refresh cleanup errors occurred",
			"operation", "refresh_token_pair_cleanup",
			"cleanup_errors", refreshErrs,
			"error_count", len(refreshErrs),
			"result", "core_operation_succeeded")
	}

	return nil
}

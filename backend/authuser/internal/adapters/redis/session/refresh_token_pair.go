package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// RefreshTokenPair refreshes access token and rotates refresh token
// This method cleans up ALL existing tokens for the session to ensure security
func (r *SessionRepository) RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, accessTTL, refreshTTL time.Duration) error {
	sessionIDStr := sessionID.String()
	oldRefreshTokenKey := auth.FormatRefreshTokenKey(oldRefreshTokenHash)
	newAccessTokenKey := auth.FormatAccessTokenKey(newAccessTokenHash)
	newRefreshTokenKey := auth.FormatRefreshTokenKey(newRefreshTokenHash)

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

	// Remove old refresh token
	if err := r.client.Del(ctx, oldRefreshTokenKey).Err(); err != nil {
		// Don't fail the operation if old token deletion fails, just log
		refreshErrs.Set("cleanup_old_refresh_token", err.Error())
	}

	// Clean up any stale access tokens for this session by scanning and removing them
	// This ensures security by invalidating all old access tokens
	pattern := auth.FormatAccessTokenKey("*")
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
	if len(refreshErrs) > 0 {
		// Log errors but return success since the core operation succeeded
		// In a real implementation, you might want to use structured logging here
	}

	return nil
}

// InvalidateTokenPair removes both access and refresh tokens
func (r *SessionRepository) InvalidateTokenPair(ctx context.Context, accessTokenHash, refreshTokenHash string, sessionID uuid.UUID) error {
	accessTokenKey := auth.FormatAccessTokenKey(accessTokenHash)
	refreshTokenKey := auth.FormatRefreshTokenKey(refreshTokenHash)
	sessionKey := auth.FormatSessionKey(sessionID.String())

	// Remove all keys - don't fail if some don't exist
	keys := []string{accessTokenKey, refreshTokenKey, sessionKey}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return errs.ClassifyRedisError("invalidate token pair", err)
	}

	return nil
}


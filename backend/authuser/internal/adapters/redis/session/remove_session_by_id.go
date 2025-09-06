package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (r *SessionRepository) RemoveSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	sessionKey := session.FormatSessionKey(sessionID.String())
	sessionIDStr := sessionID.String()

	// First, get the session data to extract token hashes and userIDHash
	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Session doesn't exist - this is not an error for removal
			return nil
		}
		return errs.ClassifyRedisError("get session data for removal", err)
	}

	// Decode session to get token hashes and userIDHash
	decodedSession, err := session.DecodeSession([]byte(sessionData))
	if err != nil {
		// If we can't decode, we can still remove the session but can't clean up tokens/index
		// This is acceptable since tokens will expire naturally
		if delErr := r.client.Del(ctx, sessionKey).Err(); delErr != nil {
			return errs.ClassifyRedisError("remove session by ID (decode failed)", delErr)
		}
		return nil
	}

	// Use pipeline for atomic operations
	pipe := r.client.Pipeline()

	// Remove session data
	pipe.Del(ctx, sessionKey)

	// Remove token mappings
	accessTokenKey := session.FormatAccessTokenKey(decodedSession.AccessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(decodedSession.RefreshTokenHash)
	pipe.Del(ctx, accessTokenKey)
	pipe.Del(ctx, refreshTokenKey)

	// Remove session ID from user session index
	userSessionIndexKey := session.FormatUserSessionIndexKey(decodedSession.UserIDHash)
	pipe.SRem(ctx, userSessionIndexKey, sessionIDStr)

	// Execute all operations
	_, err = pipe.Exec(ctx)
	if err != nil {
		return errs.ClassifyRedisError("execute session removal", err)
	}

	return nil
}

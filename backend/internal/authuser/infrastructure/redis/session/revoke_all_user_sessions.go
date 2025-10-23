package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/redis/go-redis/v9"
)

func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, userIDHash string) error {
	userSessionIndexKey := session.FormatUserSessionIndexKey(userIDHash)

	// Get all session IDs for this user
	sessionIDs, err := r.client.SMembers(ctx, userSessionIndexKey).Result()
	if err != nil {
		return errs.ClassifyRedisError("get user session IDs", err)
	}

	if len(sessionIDs) == 0 {
		// No sessions to revoke
		return nil
	}

	// Use pipeline for efficient bulk operations
	pipe := r.client.Pipeline()

	// Remove all session data and token mappings
	for _, sessionID := range sessionIDs {
		sessionKey := session.FormatSessionKey(sessionID)

		// Get the session data to extract token hashes
		sessionData, err := r.client.Get(ctx, sessionKey).Result()
		if err != nil && err != redis.Nil {
			return errs.ClassifyRedisError("get session data for cleanup", err)
		}

		// Delete session key
		pipe.Del(ctx, sessionKey)

		// If we have session data, decode it and delete token keys
		if err != redis.Nil {
			decodedSession, decodeErr := session.DecodeSession([]byte(sessionData))
			if decodeErr == nil {
				accessTokenKey := session.FormatAccessTokenKey(decodedSession.AccessTokenHash)
				refreshTokenKey := session.FormatRefreshTokenKey(decodedSession.RefreshTokenHash)
				pipe.Del(ctx, accessTokenKey)
				pipe.Del(ctx, refreshTokenKey)
			}
			// If decode fails, we still delete the session but can't clean up tokens
			// This is acceptable since tokens will expire naturally
		}
	}

	// Remove the user session index
	pipe.Del(ctx, userSessionIndexKey)

	// Execute all operations
	_, err = pipe.Exec(ctx)
	if err != nil {
		return errs.ClassifyRedisError("execute session revocation", err)
	}

	return nil
}

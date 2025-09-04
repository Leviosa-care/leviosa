package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
)

// FindSessionByRefreshToken implements two-step security: refresh token -> session ID -> session data
func (r *SessionRepository) FindSessionByRefreshToken(ctx context.Context, refreshTokenHash string) (string, []byte, error) {
	// Step 1: Get session ID from refresh token hash
	refreshTokenKey := FormatRefreshTokenKey(refreshTokenHash)
	sessionID, err := r.client.Get(ctx, refreshTokenKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session ID by refresh token", err)
	}

	// Step 2: Get session data using session ID
	sessionKey := FormatSessionKey(sessionID)
	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session data by session ID", err)
	}

	return sessionID, []byte(sessionData), nil
}


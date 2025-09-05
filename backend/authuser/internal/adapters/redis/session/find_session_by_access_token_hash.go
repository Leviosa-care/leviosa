package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
)

// FindSessionByAccessTokenHash implements two-step security: access token -> session ID -> session data
func (r *SessionRepository) FindSessionByAccessTokenHash(ctx context.Context, accessTokenHash string) (string, []byte, error) {
	// Step 1: Get session ID from access token hash
	accessTokenKey := FormatAccessTokenKey(accessTokenHash)
	sessionID, err := r.client.Get(ctx, accessTokenKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session ID by access token", err)
	}

	// Step 2: Get session data using session ID
	sessionKey := FormatSessionKey(sessionID)
	sessionData, err := r.client.Get(ctx, sessionKey).Bytes()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session data by session ID", err)
	}

	return sessionID, sessionData, nil
}

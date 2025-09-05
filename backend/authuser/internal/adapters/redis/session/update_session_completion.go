package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/google/uuid"
)

func (r *SessionRepository) UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, sessionEncoded []byte) error {
	sessionKey := auth.FormatSessionKey(sessionID.String())
	
	// Get the current TTL to preserve it
	ttl, err := r.client.TTL(ctx, sessionKey).Result()
	if err != nil {
		return errs.ClassifyRedisError("failed to get session TTL", err)
	}
	
	// If key doesn't exist, TTL returns -2
	if ttl == -2*time.Second {
		return errs.ErrRepositoryNotFound
	}
	
	// Update the session data while preserving TTL
	err = r.client.Set(ctx, sessionKey, sessionEncoded, ttl).Err()
	if err != nil {
		return errs.ClassifyRedisError("failed to update session completion", err)
	}
	
	return nil
}
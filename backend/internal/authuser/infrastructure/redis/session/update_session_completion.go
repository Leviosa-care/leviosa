package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *SessionRepository) UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, sessionEncoded []byte) error {
	sessionKey := session.FormatSessionKey(sessionID.String())

	// Get the current TTL to preserve it
	ttl, err := r.client.TTL(ctx, sessionKey).Result()
	if err != nil {
		return errs.ClassifyRedisError("failed to get session TTL", err)
	}

	// If key doesn't exist, TTL returns -2 nanoseconds (not -2 seconds)
	// Redis returns -2 for non-existent keys, go-redis returns this as time.Duration(-2)
	if ttl == time.Duration(-2) {
		return errs.ErrRepositoryNotFound
	}

	// Update the session data while preserving TTL
	err = r.client.Set(ctx, sessionKey, sessionEncoded, ttl).Err()
	if err != nil {
		return errs.ClassifyRedisError("failed to update session completion", err)
	}

	return nil
}

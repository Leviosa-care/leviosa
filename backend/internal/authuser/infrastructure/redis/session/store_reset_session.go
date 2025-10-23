package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *SessionRepository) StoreResetSession(ctx context.Context, tokenHash, userEmail string, ttl time.Duration) error {
	key := FormatResetSessionKey(tokenHash)

	if err := r.client.Set(ctx, key, userEmail, ttl).Err(); err != nil {
		return errs.ClassifyRedisError("store reset session", err)
	}

	return nil
}


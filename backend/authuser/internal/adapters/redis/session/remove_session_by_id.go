package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (r *SessionRepository) RemoveSessionByID(ctx context.Context, sessionID string) error {
	key := auth.FormatSessionKey(sessionID)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errs.ClassifyRedisError("remove session by ID", err)
	}

	return nil
}

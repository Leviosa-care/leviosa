package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (r *SessionRepository) RemoveSessionByToken(ctx context.Context, tokenHash string) error {
	key := auth.FormatTokenKey(tokenHash)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errs.ClassifyRedisError("remove session by token hash", err)
	}

	return nil
}

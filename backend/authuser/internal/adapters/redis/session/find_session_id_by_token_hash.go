package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
)

func (r *SessionRepository) FindSessionIDByTokenHash(ctx context.Context, tokenHash string) (string, error) {
	key := FormatTokenKey(tokenHash)

	sessionID, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", errs.ClassifyRedisError("get session by token hash", err)
	}

	return sessionID, nil
}

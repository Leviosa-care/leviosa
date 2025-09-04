package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
)

func (r *SessionRepository) FindSessionByID(ctx context.Context, sessionID string) ([]byte, error) {
	key := FormatSessionKey(sessionID)

	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errs.ClassifyRedisError("get session by ID", err)
	}
	return result, nil
}

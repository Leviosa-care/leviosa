package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (r *SessionRepository) FindSessionByID(ctx context.Context, sessionID string) ([]byte, error) {
	key := auth.FormatSessionKey(sessionID)

	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errs.ClassifyRedisError("get session by ID", err)
	}
	return result, nil
}

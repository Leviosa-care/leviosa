package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *SessionRepository) FindSessionByID(ctx context.Context, sessionID uuid.UUID) ([]byte, error) {
	key := session.FormatSessionKey(sessionID.String())

	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errs.ClassifyRedisError("get session by ID", err)
	}
	return result, nil
}

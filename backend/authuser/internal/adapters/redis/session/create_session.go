package sessionRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

func (r *SessionRepository) CreateSession(ctx context.Context, sessionID uuid.UUID, tokenHash string, sessionEncoded []byte, ttl time.Duration) error {
	sessionKey := FormatSessionKey(sessionID.String())
	tokenKey := FormatTokenKey(tokenHash)

	if err := r.client.Set(ctx, sessionKey, sessionEncoded, ttl).Err(); err != nil {
		return errs.ClassifyRedisError("save session with ID key", err)
	}

	var errs errsx.Map
	if err := r.client.Set(ctx, tokenKey, sessionID.String(), ttl).Err(); err != nil {
		if delErr := r.client.Del(ctx, sessionKey).Err(); delErr != nil {
			errs.Set("rollback session key", delErr.Error())
		}
		errs.Set("save token key", err.Error())
		return errs.AsError()
	}
	return nil
}

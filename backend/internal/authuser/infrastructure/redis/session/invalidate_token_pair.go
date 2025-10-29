package sessionRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// InvalidateTokenPair removes session, access and refresh tokens
func (r *SessionRepository) InvalidateTokenPair(ctx context.Context, accessTokenHash, refreshTokenHash string, sessionID uuid.UUID) error {
	accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
	refreshTokenKey := session.FormatRefreshTokenKey(refreshTokenHash)
	sessionKey := session.FormatSessionKey(sessionID.String())

	// Remove all keys - don't fail if some don't exist
	keys := []string{accessTokenKey, refreshTokenKey, sessionKey}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return errs.ClassifyRedisError("invalidate token pair", err)
	}

	return nil
}

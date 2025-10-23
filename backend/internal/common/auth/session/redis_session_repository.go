package session

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/redis/go-redis/v9"
)

const (
	SessionKeyPrefix          = "auth:session:"
	AccessTokenKeyPrefix      = "auth:access:"
	RefreshTokenKeyPrefix     = "auth:refresh:"
	UserSessionIndexKeyPrefix = "auth:user_sessions:"
)

// RedisSessionRepository implements the minimal SessionRepository interface for authentication middleware
type RedisSessionRepository struct {
	client *redis.Client
}

// NewRedisSessionRepository creates a new Redis-based session repository for authentication
func NewRedisSessionRepository(client *redis.Client) SessionRepository {
	return &RedisSessionRepository{
		client: client,
	}
}

// FindSessionByAccessTokenHash implements two-step security: access token -> session ID -> session data
func (r *RedisSessionRepository) FindSessionByAccessTokenHash(ctx context.Context, accessTokenHash string) (string, []byte, error) {
	// Step 1: Get session ID from access token hash
	accessTokenKey := FormatAccessTokenKey(accessTokenHash)
	sessionID, err := r.client.Get(ctx, accessTokenKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session ID by access token", err)
	}

	// Step 2: Get session data using session ID
	sessionKey := FormatSessionKey(sessionID)
	sessionData, err := r.client.Get(ctx, sessionKey).Bytes()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session data by session ID", err)
	}

	return sessionID, sessionData, nil
}

// FindSessionByRefreshTokenHash implements two-step security: refresh token -> session ID -> session data
func (r *RedisSessionRepository) FindSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, []byte, error) {
	// Step 1: Get session ID from refresh token hash
	refreshTokenKey := FormatRefreshTokenKey(refreshTokenHash)
	sessionID, err := r.client.Get(ctx, refreshTokenKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session ID by refresh token", err)
	}

	// Step 2: Get session data using session ID
	sessionKey := FormatSessionKey(sessionID)
	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session data by session ID", err)
	}

	return sessionID, []byte(sessionData), nil
}

// Redis key formatting functions - exported for test helpers
func FormatSessionKey(sessionID string) string {
	return fmt.Sprintf("%s%s", SessionKeyPrefix, sessionID)
}

func FormatAccessTokenKey(accessTokenHash string) string {
	return fmt.Sprintf("%s%s", AccessTokenKeyPrefix, accessTokenHash)
}

func FormatRefreshTokenKey(refreshTokenHash string) string {
	return fmt.Sprintf("%s%s", RefreshTokenKeyPrefix, refreshTokenHash)
}

func FormatUserSessionIndexKey(userIDHash string) string {
	return fmt.Sprintf("%s%s", UserSessionIndexKeyPrefix, userIDHash)
}

package auth

import "context"

// SessionRepository defines the minimal interface needed for authentication middleware
// This interface includes session retrieval methods needed for auth validation
type SessionRepository interface {
	FindSessionByAccessTokenHash(ctx context.Context, accessTokenHash string) (string, []byte, error)
	FindSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, []byte, error)
}

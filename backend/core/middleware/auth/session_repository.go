package auth

import "context"

// SessionRepository defines the minimal interface needed for authentication middleware
// This interface includes session retrieval methods needed for auth validation
type SessionRepository interface {
	// Legacy single-token authentication
	FindSessionByTokenHash(ctx context.Context, tokenHash string) ([]byte, error)

	// New dual-token authentication
	FindSessionByAccessTokenHash(ctx context.Context, accessTokenHash string) (string, []byte, error)
	FindSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, []byte, error)
}

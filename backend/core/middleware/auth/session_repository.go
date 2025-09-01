package auth

import "context"

// SessionRepository defines the minimal interface needed for authentication middleware
// This interface includes session retrieval methods needed for auth validation
type SessionRepository interface {
	// Legacy single-token authentication
	FindSessionByTokenHash(ctx context.Context, tokenHash string) ([]byte, error)

	// New dual-token authentication
	FindSessionByAccessToken(ctx context.Context, accessTokenHash string) ([]byte, error)
	FindSessionByRefreshToken(ctx context.Context, refreshTokenHash string) ([]byte, error)
}

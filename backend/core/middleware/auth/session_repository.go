package auth

import "context"

// SessionRepository defines the minimal interface needed for authentication middleware
// This interface only includes session retrieval needed for auth validation
type SessionRepository interface {
	FindSessionByTokenHash(ctx context.Context, tokenHash string) ([]byte, error)
}

package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/google/uuid"
)

// SessionRepository defines the complete interface for session management
// It embeds the minimal auth.SessionRepository interface and adds additional operations
type SessionRepository interface {
	// Embed the minimal authentication interface from core
	auth.SessionRepository

	// Additional operations needed for full session management
	FindSessionByID(ctx context.Context, sessionID string) ([]byte, error)
	FindSessionIDByTokenHash(ctx context.Context, tokenHash string) (string, error)
	CreateSession(ctx context.Context, sessionID uuid.UUID, tokenHash string, sessionEncoded []byte, ttl time.Duration) error
	RemoveSessionByID(ctx context.Context, sessionID string) error
	RemoveSessionByToken(ctx context.Context, sessionID string) error
}
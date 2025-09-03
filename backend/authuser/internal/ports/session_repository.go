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

	// Legacy single-token operations
	FindSessionByID(ctx context.Context, sessionID string) ([]byte, error)
	FindSessionIDByTokenHash(ctx context.Context, tokenHash string) (string, error)
	CreateSession(ctx context.Context, sessionID uuid.UUID, tokenHash string, sessionEncoded []byte, ttl time.Duration) error
	RemoveSessionByID(ctx context.Context, sessionID string) error
	RemoveSessionByToken(ctx context.Context, sessionID string) error

	// New dual-token operations
	FindSessionByAccessToken(ctx context.Context, accessTokenHash string) (string, []byte, error)
	FindSessionByRefreshToken(ctx context.Context, refreshTokenHash string) (string, []byte, error)
	CreateTokenPair(ctx context.Context, sessionID uuid.UUID, accessTokenHash, refreshTokenHash string, sessionEncoded []byte, accessTTL, refreshTTL time.Duration) error
	RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, accessTTL, refreshTTL time.Duration) error
	InvalidateTokenPair(ctx context.Context, accessTokenHash, refreshTokenHash string, sessionID uuid.UUID) error

	// Session state operations
	UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, sessionEncoded []byte) error
}

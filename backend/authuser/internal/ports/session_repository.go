package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionRepository defines the complete interface for session management
type SessionRepository interface {
	CreateTokenPair(ctx context.Context, sessionID uuid.UUID, accessTokenHash, refreshTokenHash string, sessionEncoded []byte, accessTTL, refreshTTL time.Duration) error
	FindSessionByID(ctx context.Context, sessionID uuid.UUID) ([]byte, error)
	RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, accessTTL, refreshTTL time.Duration) error
	RevokeAllUserSessions(ctx context.Context, UserID uuid.UUID) error
	UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, sessionEncoded []byte) error
}

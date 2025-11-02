package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionRepository defines the complete interface for session management
type SessionRepository interface {
	CreateSession(ctx context.Context, sessionID uuid.UUID, accessTokenHash, refreshTokenHash, userIDHash string, sessionEncoded []byte, accessTTL, refreshTTL time.Duration) error
	FindSessionByID(ctx context.Context, sessionID uuid.UUID) ([]byte, error)
	RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, updatedSessionData []byte, accessTTL, refreshTTL time.Duration) error
	RevokeAllUserSessions(ctx context.Context, userIDHash string) error
	UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, sessionEncoded []byte) error
	RemoveSessionByID(ctx context.Context, sessionID uuid.UUID) error
	StoreResetSession(ctx context.Context, tokenHash, userEmail string, ttl time.Duration) error
	ValidateResetSession(ctx context.Context, tokenHash string) (string, error)
	InvalidateTokenPair(ctx context.Context, accessTokenHash, refreshTokenHash string, sessionID uuid.UUID) error
}

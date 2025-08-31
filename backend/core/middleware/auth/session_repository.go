package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	FindSessionByID(ctx context.Context, sessionID string) ([]byte, error)
	FindSessionIDByTokenHash(ctx context.Context, tokenHash string) (string, error)
	FindSessionByTokenHash(ctx context.Context, tokenHash string) ([]byte, error)
	CreateSession(ctx context.Context, sessionID uuid.UUID, tokenHash string, sessionEncoded []byte, ttl time.Duration) error
	RemoveSessionByID(ctx context.Context, sessionID string) error
	RemoveSessionByToken(ctx context.Context, sessionID string) error
}
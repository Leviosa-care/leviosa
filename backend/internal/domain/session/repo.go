package sessionService

import (
	"context"
)

type SessionRepository interface {
	FindSessionByID(ctx context.Context, sessionID string) ([]byte, error)
	CreateSession(ctx context.Context, sessionID string, sessionEncoded []byte) error
	RemoveSession(ctx context.Context, sessionID string) error
}

type Reader interface {
	FindSessionByID(ctx context.Context, sessionID string) ([]byte, error)
}

type Writer interface {
	CreateSession(ctx context.Context, sessionID string, sessionEncoded []byte) error
	RemoveSession(ctx context.Context, sessionID string) error
}

type ReadWriter interface {
	Reader
	Writer
}

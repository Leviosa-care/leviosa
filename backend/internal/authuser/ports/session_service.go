package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(ctx context.Context, request *domain.CreateSessionRequest) (*domain.CreateSessionResponse, error)
	// IsSessionValid(ctx context.Context, request *domain.ValidateSessionRequest) (bool, error)
	RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error)
	RemoveSession(ctx context.Context, sessionID uuid.UUID) error
	UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, completedAt *time.Time) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	CreateResetSession(ctx context.Context, token, userEmail string, ttl time.Duration) error
	ValidateResetSession(ctx context.Context, token string) (string, error)
}

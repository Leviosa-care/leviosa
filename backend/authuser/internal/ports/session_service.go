package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(ctx context.Context, request *domain.CreateSessionRequest) (*domain.CreateSessionResponse, error)
	// IsSessionValid(ctx context.Context, request *domain.ValidateSessionRequest) (bool, error)
	RefreshSession(ctx context.Context, sessionID uuid.UUID) (*domain.RefreshSessionResponse, error)
	UpdateSessionCompletion(ctx context.Context, sessionID uuid.UUID, completedAt *time.Time) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
}

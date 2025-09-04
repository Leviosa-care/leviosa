package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/middleware/auth"
)

type SessionService interface {
	CreateSession(ctx context.Context, request *domain.CreateSessionRequest) (*domain.CreateSessionResponse, error)
	GetSession(ctx context.Context, request *domain.GetSessionRequest) (*auth.Session, error)
	RemoveSession(ctx context.Context, request *domain.RemoveSessionRequest) error
	IsSessionValid(ctx context.Context, request *domain.ValidateSessionRequest) (bool, error)
	RefreshSession(ctx context.Context, request *domain.RefreshSessionRequest) (*domain.RefreshSessionResponse, error)
	UpdateSessionCompletion(ctx context.Context, sessionID string, completedAt *time.Time) error
}

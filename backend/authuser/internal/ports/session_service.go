package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/middleware/auth"
)

type SessionService interface {
	CreateSession(ctx context.Context, request *domain.CreateSessionRequest) (*domain.CreateSessionResponse, error)
}

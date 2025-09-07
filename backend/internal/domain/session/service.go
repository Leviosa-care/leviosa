package sessionService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

// type Service interface {
type SessionService interface {
	CreateSession(ctx context.Context, userID string, role models.Role) (string, error)
	RemoveSession(ctx context.Context, sessionID string) error
}

type service struct {
	Repo ReadWriter
}

func New(repo ReadWriter) Service {
	return &service{
		Repo: repo,
	}
}

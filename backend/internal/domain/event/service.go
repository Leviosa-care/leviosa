package eventService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	"github.com/hengadev/leviosa/internal/domain/event/security"
	"github.com/hengadev/leviosa/pkg/config"
)

type Service interface {
	CreateEvent(ctx context.Context, event *models.Event) (string, error)
	DecreasePlacecount(ctx context.Context, eventID string) error
	RemoveEvent(ctx context.Context, eventID string) error
	ModifyEvent(ctx context.Context, event *models.Event) error
}

type service struct {
	repo ReadWriter
	*security.SecureEventData
}

func New(repo ReadWriter, config *config.SecurityConfig) Service {
	return &service{
		repo,
		security.NewSecureEventData(config),
	}
}

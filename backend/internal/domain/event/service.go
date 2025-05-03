package eventService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/event/models"

	"github.com/hengadev/encx"
)

type Service interface {
	CreateEvent(ctx context.Context, event *models.Event) (string, error)
	DecreasePlacecount(ctx context.Context, eventID string) error
	RemoveEvent(ctx context.Context, eventID string) error
	ModifyEvent(ctx context.Context, event *models.Event) error
}

type service struct {
	repo   ReadWriter
	crypto *encx.Crypto
}

func New(repo ReadWriter, crypto *encx.Crypto) Service {
	return &service{
		repo,
		crypto,
	}
}

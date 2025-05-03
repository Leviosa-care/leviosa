package registerService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/event/models"
)

type Service interface {
	AddRegistration(ctx context.Context, registration *Registration) error
	CreateRegistration(ctx context.Context, userID, spotStr string, event *models.Event) error
}
type service struct {
	Repo ReadWriter
}

func NewService(repo ReadWriter) Service {
	return &service{Repo: repo}
}

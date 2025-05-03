package throttlerService

import "context"

type Service interface {
	RegisterAttempt(ctx context.Context, email string) error
	Reset(ctx context.Context, email string) error
}

type service struct {
	repo ReadWriter
}

func New(repo ReadWriter) Service {
	return &service{
		repo: repo,
	}
}

package session

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/hengadev/encx"
)

type SessionService struct {
	repo   ports.SessionRepository
	crypto encx.CryptoService
}

func New(ctx context.Context, repo ports.SessionRepository, crypto encx.CryptoService) ports.SessionService {
	return &SessionService{repo: repo, crypto: crypto}
}

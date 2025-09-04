package session

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/hengadev/encx"
)

type SessionService struct {
	repo   ports.SessionRepository
	crypto encx.CryptoService
	cache  *TokenDurationCache
}

func New(ctx context.Context, repo ports.SessionRepository, crypto encx.CryptoService) ports.SessionService {
	return &SessionService{
		repo:   repo,
		crypto: crypto,
		cache:  NewTokenDurationCache(),
	}
}

func NewWithCache(ctx context.Context, repo ports.SessionRepository, crypto encx.CryptoService, cache *TokenDurationCache) ports.SessionService {
	return &SessionService{
		repo:   repo,
		crypto: crypto,
		cache:  cache,
	}
}

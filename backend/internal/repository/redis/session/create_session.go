package sessionRepository

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/session"
	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/redis/go-redis/v9"
)

func (s *Repository) CreateSession(ctx context.Context, sessionID string, sessionEncoded []byte) error {
	key := formatSessionKey(sessionID)
	result := s.client.Set(ctx, key, sessionEncoded, sessionService.SessionDuration)
	if err := result.Err(); err != nil {
		switch {
		case errors.Is(err, redis.ErrClosed):
			return rp.NewDatabaseErr(err)
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}

	if result.Val() == "" {
		return rp.NewNotCreatedErr(nil, "session")
	}
	return nil
}

package sessionRepository

import (
	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func New(client *redis.Client) ports.SessionRepository {
	return &SessionRepository{client: client}
}

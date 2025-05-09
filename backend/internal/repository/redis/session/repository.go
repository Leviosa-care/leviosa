package sessionRepository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const SESSIONPREFIX = "session"
const sessionKeyFormat = "session:%s"

type Repository struct {
	client *redis.Client
}

func New(ctx context.Context, client *redis.Client) *Repository {
	return &Repository{client}
}

func (r *Repository) GetClient() *redis.Client {
	return r.client
}

func formatSessionKey(emailHash string) string {
	return fmt.Sprintf(sessionKeyFormat, emailHash)
}

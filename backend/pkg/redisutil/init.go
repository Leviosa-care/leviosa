package redisutil

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/session"
	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/redis/go-redis/v9"
)

// TODO: find a value for the expiration of the value set
func Init(ctx context.Context, client *redis.Client, queries map[string]interface{}) error {
	for k, v := range queries {
		err := client.Set(ctx, k, v, sessionService.SessionDuration).Err()
		if err != nil {
			return rp.NewNotCreatedErr(err, "query")
		}
	}
	return nil
}

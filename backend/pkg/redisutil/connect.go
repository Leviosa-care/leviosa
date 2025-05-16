package redisutil

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, opts ...RedisOption) (*redis.Client, error) {
	r := &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
	for _, opt := range opts {
		opt(r)
	}
	client := redis.NewClient(r)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return client, nil
}

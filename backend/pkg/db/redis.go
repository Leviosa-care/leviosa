package db

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/session"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/pkg/config"

	"github.com/redis/go-redis/v9"
)

type RedisOption func(*redis.Options)

func WithDB(DB int) RedisOption {
	return func(r *redis.Options) {
		r.DB = DB
	}
}

func WithAddr(addr string) RedisOption {
	return func(r *redis.Options) {
		r.Addr = addr
	}
}

func WithPassword(pwd string) RedisOption {
	return func(r *redis.Options) {
		r.Password = pwd
	}
}

func DefaultRedis() *redis.Options {
	return &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
}

func Redis(ctx context.Context, options *config.RedisSecrets) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return client, nil
}

func Init(ctx context.Context, client *redis.Client, queries map[string]any) error {
	for k, v := range queries {
		err := client.Set(ctx, k, v, sessionService.SessionDuration).Err()
		if err != nil {
			return rp.NewNotCreatedErr(err, "initial query")
		}
	}
	return nil
}

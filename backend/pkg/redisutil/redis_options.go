package redisutil

import "github.com/redis/go-redis/v9"

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

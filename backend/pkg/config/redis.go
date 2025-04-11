package config

import (
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/pkg/flags"

	"github.com/go-redis/redis"
)

type redisCreds struct {
	*redis.Options
}

func (c *Config) GetRedis() *redisCreds {
	return c.redis
}

func (c *Config) setRedis(env mode.EnvMode) error {
	var addr, password string
	var db int
	switch env {
	case mode.ModeDev:
		addr = "127.0.0.1:6379"
		password = "secret"
		db = 0
	case mode.ModeStaging, mode.ModeProd:
		addr = c.viper.GetString("redis.addr")
		password = c.viper.GetString("redis.password")
		db = c.viper.GetInt("redis.db")
	default:
		return fmt.Errorf("mode value can only be 'development', 'production' or 'staging', got : %q", env)
	}
	if addr == "" {
		return errors.New("'REDIS_ADDR' environment variable not set; please define it to specify Redis address")
	}
	if password == "" {
		return errors.New("'REDIS_PASSWORD' environment variable not set; please define it to specify Redis password")
	}
	if db >= 16 || db < 0 {
		return errors.New("'REDIS_DB' environment variable not set; please define it to specify Redis database")
	}
	c.redis = &redisCreds{
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		},
	}
	return nil
}

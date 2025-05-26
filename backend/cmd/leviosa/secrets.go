package main

import (
	"context"

	"github.com/hengadev/leviosa/pkg/config"

	"github.com/hengadev/errsx"
)

const path = "/secrets/leviosa-secrets.json"

func loadSecrets(ctx context.Context) (*config.RedisSecrets, *config.PostgresSecrets, *config.RabbitSecrets, error) {
	manager := config.NewManager(config.DefaultConfig())
	redisLoader := config.RegisterLoader[config.RedisSecrets](manager, "redis", nil, opts.mode)
	postgresLoader := config.RegisterLoader[config.PostgresSecrets](manager, "postgres", nil, opts.mode)
	rabbitLoader := config.RegisterLoader[config.RabbitSecrets](manager, "rabbit", nil, opts.mode)
	var errs errsx.Map
	redisConf, err := redisLoader.Load(ctx, path)
	if err != nil {
		errs.Set(redisConf.GetType(), err)
	}
	postgresConf, err := postgresLoader.Load(ctx, path)
	if err != nil {
		errs.Set(postgresConf.GetType(), err)
	}
	rabbitConf, err := rabbitLoader.Load(ctx, path)
	if err != nil {
		errs.Set(rabbitConf.GetType(), err)
	}
	return &redisConf, &postgresConf, &rabbitConf, errs.AsError()
}

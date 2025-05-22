package main

import (
	"context"
	"database/sql"
	"fmt"

	cfg "github.com/hengadev/leviosa/internal/config"
	"github.com/hengadev/leviosa/pkg/db"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
)

func setupDatabases(
	ctx context.Context,
	conf *cfg.Config,
	env envmode.Mode,
) (
	*sql.DB,
	*redis.Client,
	*s3.Client,
	error,
) {
	redisConfig := conf.GetRedis()
	redisClient, err := db.Redis(ctx, redisConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating connection to redis database: %w", err)
	}
	db, err := db.SQLite(ctx, env)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating connection to sqlite database: %w", err)
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load default configuration for S3 repository: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	return db, redisClient, s3Client, nil
}

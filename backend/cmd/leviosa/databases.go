package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/leviosa/pkg/db"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/redis/go-redis/v9"
)

func setupDatabases(
	ctx context.Context,
	redisOptions *redis.Options,
	env envmode.Mode,
) (*sql.DB, *redis.Client, error) {
	client, err := db.Redis(ctx, redisOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("creating connection to redis database: %w", err)
	}
	db, err := db.SQLite(ctx, env)
	if err != nil {
		return nil, nil, fmt.Errorf("creating connection to sqlite database: %w", err)
	}
	return db, client, nil
}

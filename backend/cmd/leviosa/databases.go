package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/leviosa/internal/config"
	"github.com/hengadev/leviosa/pkg/envmode"
	"github.com/hengadev/leviosa/pkg/redisutil"
	"github.com/hengadev/leviosa/pkg/sqliteutil"

	"github.com/redis/go-redis/v9"
)

func setupDatabases(
	ctx context.Context,
	conf *config.Config,
	env envmode.Mode,
) (*sql.DB, *redis.Client, error) {
	sqliteConf := conf.GetSQLITE()
	redisConf := conf.GetRedis()

	// databases setup
	redisdb, err := redisutil.Connect(
		ctx,
		redisutil.WithAddr(redisConf.Addr),
		redisutil.WithDB(redisConf.DB),
		redisutil.WithPassword(redisConf.Password),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating connection to redis database: %w", err)
	}

	sqliteDSN := sqliteutil.BuildDSN(env, sqliteConf.Filename)
	sqlitedb, err := sqliteutil.Connect(ctx, sqliteDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("creating connection to sqlite database: %w", err)
	}
	return sqlitedb, redisdb, nil
}

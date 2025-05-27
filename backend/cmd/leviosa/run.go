package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hengadev/leviosa/internal/config"
	"github.com/hengadev/leviosa/internal/server"
	"github.com/hengadev/leviosa/internal/server/app"
	// "github.com/hengadev/leviosa/internal/server/cron"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/joho/godotenv"
)

var opts struct {
	mode   envmode.Mode
	server struct {
		port int
	}
	logger struct {
		style string
		level string
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// always load the .env file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load env variables: %w", err)
	}

	// set environment file (using [mode].env for specified mode)
	if opts.mode == envmode.Dev {
		if err := godotenv.Load(fmt.Sprintf("%s.env", opts.mode.String())); err != nil {
			return fmt.Errorf("loading env variables: %w", err)
		}
	}

	// setup env variables
	if err := setupEnvVars(); err != nil {
		return fmt.Errorf("failed to get env variables: %w", err)
	}

	// set logger
	slogHandler, err := setLogger()
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	// NOTE: with the new implementation of the config
	// configuration and secret handling
	// redisConf, postgresConf, rabbitConf, err := loadSecrets(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to load secrets: %w", err)
	// }
	// NOTE: old version
	cfg, err := config.Load(ctx, opts.mode)
	if err != nil {
		return fmt.Errorf("failed to load secrets: %w", err)
	}
	redisConf, postgresConf, rabbitConf, bucketname := cfg.GetRedis(), cfg.GetPostgres(), cfg.GetRabbitMQ(), cfg.GetS3().BucketName

	postgresdb, redisdb, s3Client, err := setupDatabases(ctx, redisConf, postgresConf, opts.mode)
	if err != nil {
		return fmt.Errorf("setting up databases: %w", err)
	}

	rabbitConn, err := setBroker(ctx, rabbitConf)
	defer rabbitConn.Close()

	appSvcs, appRepos, err := makeServices(ctx, postgresdb, redisdb, s3Client, rabbitConn, bucketname)
	if err != nil {
		return fmt.Errorf("create services: %w", err)
	}
	appCtx := app.New(&appSvcs, &appRepos)
	srv := server.New(
		appCtx,
		opts.mode,
		slogHandler,
		server.WithPort(opts.server.port),
	)
	var srvErrCh = make(chan error)

	// setting cron jobs
	// go func() {
	// 	cronHandler := cron.New(handler, logger)
	// 	if err := cronHandler.Start(); err != nil {
	// 		srvErrCh <- fmt.Errorf("cron service failed: %w", err)
	// 		return
	// 	}
	// }()

	go func() {
		slog.InfoContext(ctx, fmt.Sprintf("Server running on port %d.", opts.server.port))
		if err := srv.ListenAndServe(); err != nil {
			srvErrCh <- fmt.Errorf("launch server: %w", err)
			return
		}
	}()

	select {
	case done := <-ctx.Done():
		return fmt.Errorf("ctx.Done: %v", done)
	case err := <-srvErrCh:
		return fmt.Errorf("server error: %w", err)
	}
}

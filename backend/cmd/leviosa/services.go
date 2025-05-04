package main

import (
	"context"
	"database/sql"
	"fmt"

	// domain
	"github.com/hengadev/leviosa/internal/domain/event"
	"github.com/hengadev/leviosa/internal/domain/mail"
	"github.com/hengadev/leviosa/internal/domain/media"
	"github.com/hengadev/leviosa/internal/domain/message"
	"github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/internal/domain/product"
	"github.com/hengadev/leviosa/internal/domain/register"
	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/stripe"
	"github.com/hengadev/leviosa/internal/domain/throttler"
	"github.com/hengadev/leviosa/internal/domain/user"
	"github.com/hengadev/leviosa/internal/domain/vote"
	"github.com/hengadev/leviosa/internal/server/app"

	// repositories
	"github.com/hengadev/leviosa/internal/repository/redis/otp"
	"github.com/hengadev/leviosa/internal/repository/redis/session"
	"github.com/hengadev/leviosa/internal/repository/redis/throttler"
	"github.com/hengadev/leviosa/internal/repository/s3"
	"github.com/hengadev/leviosa/internal/repository/sqlite/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite/message"
	"github.com/hengadev/leviosa/internal/repository/sqlite/product"
	"github.com/hengadev/leviosa/internal/repository/sqlite/register"
	"github.com/hengadev/leviosa/internal/repository/sqlite/user"
	"github.com/hengadev/leviosa/internal/repository/sqlite/vote"

	// config
	"github.com/hengadev/leviosa/pkg/config"

	// external packages
	"github.com/hengadev/encx"
	"github.com/hengadev/encx/providers/hashicorpvault"
	rd "github.com/redis/go-redis/v9"
)

func makeServices(
	ctx context.Context,
	sqlitedb *sql.DB,
	redisdb *rd.Client,
	config *config.Config,
) (app.Services, app.Repos, error) {
	var appSvcs app.Services
	var appRepos app.Repos

	// crypto
	kms, err := hashicorpvault.New()
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("creating vault: %w", err)
	}
	crypto, err := encx.New(ctx, kms, "leviosa-app-key", "secret/data/pepper")
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("creating crypto: %w", err)
	}

	// user
	userRepo := userRepository.New(ctx, sqlitedb)
	userSvc := userService.New(userRepo, crypto)
	// session
	sessionRepo := sessionRepository.New(ctx, redisdb)
	sessionSvc := sessionService.New(sessionRepo)
	// event
	eventRepo := eventRepository.New(ctx, sqlitedb)
	eventSvc := eventService.New(eventRepo, crypto)
	// vote
	voteRepo := voteRepository.New(ctx, sqlitedb)
	voteSvc := vote.New(voteRepo)
	// register
	registerRepo := registerRepository.New(ctx, sqlitedb)
	registerSvc := registerService.NewService(registerRepo)
	// stripe
	stripeSvc := stripeService.New()
	// media
	bucketName := config.GetS3().BucketName
	mediaRepo, err := mediaRepository.New(ctx, bucketName)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create photo repository: %w", err)
	}
	mediaSvc := mediaService.New(mediaRepo)

	// mail
	mailSvc, err := mailService.New()
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create mail service: %w", err)
	}
	// product
	productRepo := productRepository.New(ctx, sqlitedb)
	productSvc := productService.New(productRepo)

	// throttle
	throttlerRepo := throttlerRepository.New(ctx, redisdb)
	throttlerSvc := throttlerService.New(throttlerRepo)

	// OTP
	otpRepo := otpRepository.New(ctx, redisdb)
	otpSvc := otpService.New(otpRepo)

	// message
	messageRepo := messageRepository.New(ctx, sqlitedb)
	messageSvc := messageService.New(messageRepo, crypto)

	// services
	appSvcs = app.Services{
		User:      userSvc,
		Event:     eventSvc,
		Vote:      voteSvc,
		Register:  registerSvc,
		Media:     mediaSvc,
		Session:   sessionSvc,
		Throttler: throttlerSvc,
		Mail:      mailSvc,
		Stripe:    stripeSvc,
		Product:   productSvc,
		OTP:       otpSvc,
		Message:   messageSvc,
	}
	// repos
	appRepos = app.Repos{
		User:        userRepo,
		Event:       eventRepo,
		Vote:        voteRepo,
		Register:    registerRepo,
		Media:       mediaRepo,
		Session:     sessionRepo,
		Throttler:   throttlerRepo,
		Product:     productRepo,
		OTP:         otpRepo,
		Message:     messageRepo,
		SQLiteDB:    sqlitedb,
		RedisClient: redisdb,
	}
	return appSvcs, appRepos, nil
}

package main

import (
	"context"
	"database/sql"
	"fmt"

	// domain
	"github.com/hengadev/leviosa/internal/domain/event"
	"github.com/hengadev/leviosa/internal/domain/mail"
	"github.com/hengadev/leviosa/internal/domain/message"
	"github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/internal/domain/product"
	"github.com/hengadev/leviosa/internal/domain/register"
	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/settings"
	"github.com/hengadev/leviosa/internal/domain/stripe"
	"github.com/hengadev/leviosa/internal/domain/throttler"
	"github.com/hengadev/leviosa/internal/domain/user"
	"github.com/hengadev/leviosa/internal/domain/vote"
	"github.com/hengadev/leviosa/internal/server/app"

	// repositories
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/postgres/message"
	"github.com/hengadev/leviosa/internal/repository/postgres/product"
	"github.com/hengadev/leviosa/internal/repository/postgres/register"
	"github.com/hengadev/leviosa/internal/repository/postgres/settings"
	"github.com/hengadev/leviosa/internal/repository/postgres/user"
	"github.com/hengadev/leviosa/internal/repository/postgres/vote"
	"github.com/hengadev/leviosa/internal/repository/redis/otp"
	"github.com/hengadev/leviosa/internal/repository/redis/session"
	"github.com/hengadev/leviosa/internal/repository/redis/throttler"
	"github.com/hengadev/leviosa/internal/repository/s3/settings"

	// external packages
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hengadev/encx"
	"github.com/hengadev/encx/providers/hashicorpvault"
	amqp "github.com/rabbitmq/amqp091-go"
	rd "github.com/redis/go-redis/v9"
)

const KEKAlias = "leviosa-app-key"

func makeServices(
	ctx context.Context,
	postgresdb *sql.DB,
	redisdb *rd.Client,
	s3Client *s3.Client,
	rabbitConn *amqp.Connection,
	bucketname string,
) (app.Services, app.Repos, error) {
	var appSvcs app.Services
	var appRepos app.Repos

	// crypto
	kms, err := hashicorpvault.New()
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create vault: %w", err)
	}
	crypto, err := encx.New(ctx, kms, KEKAlias, "secret/data/pepper")
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create crypto: %w", err)
	}

	// user
	userRepo, err := userRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create user repository: %w", err)
	}
	userSvc := userService.New(userRepo, crypto)
	// session
	sessionRepo := sessionRepository.New(ctx, redisdb)
	sessionSvc := sessionService.New(sessionRepo)
	// event
	eventRepo, err := eventRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create event repository: %w", err)
	}
	eventSvc := eventService.New(eventRepo, crypto)
	// vote
	voteRepo, err := voteRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create vote repository: %w", err)
	}
	voteSvc := vote.New(voteRepo)
	// register
	registerRepo, err := registerRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create register repository: %w", err)
	}
	registerSvc := registerService.NewService(registerRepo)
	// stripe
	stripeSvc := stripeService.New()
	// product
	productRepo, err := productRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create product repository: %w", err)
	}
	productSvc := productService.New(productRepo)

	// throttle
	throttlerRepo := throttlerRepository.New(ctx, redisdb)
	throttlerSvc := throttlerService.New(throttlerRepo)

	// settings
	settingsRepo, err := settingsRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create settings repository: %w", err)
	}
	settingsS3, err := settingsMedia.New(ctx, s3Client, bucketname)
	settingsSvc := settings.New(settingsRepo, settingsS3, crypto, rabbitConn)

	// OTP
	otpRepo := otpRepository.New(ctx, redisdb)
	otpSvc, err := otpService.New(ctx, crypto, otpRepo, settingsRepo, rabbitConn)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create OTP service: %w", err)
	}

	// mail
	mailSvc, err := mailService.New(ctx, settingsRepo, settingsS3, rabbitConn)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create mail service: %w", err)
	}

	// message
	messageRepo, err := messageRepository.New(ctx, postgresdb)
	if err != nil {
		return appSvcs, appRepos, fmt.Errorf("create message repository: %w", err)
	}
	messageSvc := messageService.New(messageRepo, crypto)

	// services
	appSvcs = app.Services{
		User:      userSvc,
		Event:     eventSvc,
		Vote:      voteSvc,
		Register:  registerSvc,
		Session:   sessionSvc,
		Throttler: throttlerSvc,
		Mail:      mailSvc,
		Stripe:    stripeSvc,
		Product:   productSvc,
		Settings:  settingsSvc,
		OTP:       otpSvc,
		Message:   messageSvc,
	}
	// repos
	appRepos = app.Repos{
		User:        userRepo,
		Event:       eventRepo,
		Vote:        voteRepo,
		Register:    registerRepo,
		Session:     sessionRepo,
		Throttler:   throttlerRepo,
		Product:     productRepo,
		Settings:    settingsRepo,
		OTP:         otpRepo,
		Message:     messageRepo,
		SQLiteDB:    postgresdb,
		RedisClient: redisdb,
	}
	return appSvcs, appRepos, nil
}

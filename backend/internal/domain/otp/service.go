package otpService

import (
	"context"
	"fmt"

	"github.com/hengadev/leviosa/internal/broker/rabbitmq"
	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/settings"

	"github.com/hengadev/encx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service interface {
	RequestOTP(ctx context.Context, email string) (string, error)
	VerifyOTP(ctx context.Context, email string, code string) error
	CancelOTP(ctx context.Context, email string) error
	ResendOTP(ctx context.Context, email string) (*OTP, error)
}

type service struct {
	repo     ReadWriter
	crypto   encx.CryptoService
	settings settings.Reader
	cache    *cache
}

func New(ctx context.Context, crypto encx.CryptoService, repo ReadWriter, settingsRepo settings.Reader, rabbitConn *amqp.Connection) (Service, error) {
	durationSetting, err := settingsRepo.GetInt(ctx, settings.OTPDurationKey)
	if err != nil {

	}
	duration := durationSetting.Value
	lengthSetting, err := settingsRepo.GetInt(ctx, settings.OTPLengthKey)
	if err != nil {

	}
	length := lengthSetting.Value
	maxAttemptsSetting, err := settingsRepo.GetInt(ctx, settings.OTPMaxAttemptsKey)
	if err != nil {

	}
	maxAttempts := maxAttemptsSetting.Value

	cache := newCache(duration, length, maxAttempts)
	service := &service{
		repo:     repo,
		crypto:   crypto,
		settings: settingsRepo,
		cache:    cache,
	}

	ch, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, domain.NewNotCreatedErr(fmt.Errorf("consumer channel for mail service"))
	}
	service.StartOTPSettingConsumer(ctx, ch)
	return service, nil
}

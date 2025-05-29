package mailService

import (
	"context"
	"fmt"
	"os"

	"github.com/hengadev/leviosa/internal/broker/rabbitmq"
	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/settings"
	"github.com/hengadev/leviosa/internal/domain/user/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Service interface {
	// main
	HandlePasswordForgotten(ctx context.Context, to string) error
	NewEvent(ctx context.Context, users []*models.User, eventTime string) error
	NewPayment(ctx context.Context, user *models.User, eventTime string) error
	NewVote(ctx context.Context, user *models.User, eventTime string) error
	SendRegistrationReminderEmail(ctx context.Context, user *models.User, registrationName string, daysLeft int) error
	SendOTP(ctx context.Context, email string, otp string) error
	WelcomeUser(ctx context.Context, email string) error
	// cache
	SetCompanyEmail(email string)
	SetLogo(logo []byte)
}

type service struct {
	email    string
	password string
	cache    *cache
}

func New(
	ctx context.Context,
	repo settings.Reader,
	media settings.MediaReader,
	rabbitConn *amqp.Connection,
	// ch *amqp.Channel,
) (Service, error) {
	email := os.Getenv("GMAIL_EMAIL")
	if email == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_EMAIL'"))
	}
	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_PASSWORD'"))
	}
	fromSetting, err := repo.GetString(ctx, settings.CompanyEmailKey)
	if err != nil {
		return nil, domain.NewNotFoundErr(fmt.Errorf("company email used as 'from' header"))
	}
	from := fromSetting.Value
	addressSetting, err := repo.GetString(ctx, settings.CompanyLegalAddressKey)
	if err != nil {
		return nil, domain.NewNotFoundErr(fmt.Errorf("company email used as 'address' header"))
	}
	address := addressSetting.Value

	instaSetting, err := repo.GetString(ctx, settings.CompanyEmailKey)
	if err != nil {
		return nil, domain.NewNotFoundErr(fmt.Errorf("company email used as 'insta' header"))
	}
	insta := instaSetting.Value

	// logo, err := media.GetLogo(ctx)
	// if err != nil {
	// 	return nil, domain.NewNotFoundErr(fmt.Errorf("company logo used in email templates"))
	// }
	logo := []byte{}
	cache := newCache(from, insta, address, logo)
	service := &service{
		email:    email,
		password: password,
		cache:    cache,
	}
	ch, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, domain.NewNotCreatedErr(fmt.Errorf("consumer channel for mail service"))
	}
	service.StartMailSettingConsumer(ctx, ch)
	return service, nil
}

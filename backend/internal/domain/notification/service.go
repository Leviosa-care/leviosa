package notification

import (
	"context"
	"fmt"
	"os"

	"github.com/hengadev/leviosa/internal/broker/rabbitmq"
	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/settings"
	"github.com/hengadev/leviosa/internal/domain/user/models"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Service interface {
	MailService
	SMSClient
}

type SMSClient interface {
	SendSMS(ctx context.Context, phone, message string) error
}

type smsClient struct {
	*openapi.ApiService
	sender     string
	accountSid string
}

type MailService interface {
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

type mailService struct {
	email    string
	password string
	cache    *cache
}

type service struct {
	*mailService
	*smsClient
}

// type service struct {
// 	mailService
// 	smsClient
// }

func New2(
	ctx context.Context,
	repo settings.Reader,
	media settings.MediaReader,
	rabbitConn *amqp.Connection,
	accountSid, authToken string,
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
	mailSvc := mailService{
		email:    email,
		password: password,
		cache:    cache,
	}
	ch, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, domain.NewNotCreatedErr(fmt.Errorf("consumer channel for mail service"))
	}
	mailSvc.StartMailSettingConsumer(ctx, ch)

	// sms service
	c := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	sender := os.Getenv("TWILIO_PHONE_NUMBER")
	if sender == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'TWILIO_PHONE_NUMBER'"))
	}

	return service{
		mailService: &mailSvc,
		smsClient: &smsClient{
			ApiService: c.Api,
			sender:     sender,
			accountSid: accountSid,
		},
	}, nil
}

package mailService

import (
	"context"
	"fmt"
	"os"

	"github.com/hengadev/leviosa/internal/domain"
	otpService "github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/internal/domain/user/models"
)

type Service interface {
	HandlePasswordForgotten(to string) error
	NewEvent(ctx context.Context, users []*models.User, eventTime string) error
	NewPayment(ctx context.Context, user *models.User, eventTime string) error
	NewVote(ctx context.Context, user *models.User, eventTime string) error
	PendingUser(ctx context.Context, user *models.User) error
	SendRegistrationReminderEmail(ctx context.Context, user *models.User, registrationName string, daysLeft int) error
	SendOTP(ctx context.Context, email, firstname string, otp *otpService.OTP) error
	WelcomeUser(ctx context.Context, user *models.User) error
}

type service struct {
	from     string
	email    string
	password string
}

func New() (Service, error) {
	email := os.Getenv("GMAIL_EMAIL")
	if email == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_EMAIL'"))
	}
	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_PASSWORD'"))
	}
	return &service{
		from:     "contact@leviosa.care",
		email:    email,
		password: password,
	}, nil
}

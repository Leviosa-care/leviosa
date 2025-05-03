package mailService

import (
	"context"
	"fmt"
	"os"

	"github.com/hengadev/leviosa/internal/domain"
	otpService "github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/pkg/errsx"
)

type Service interface {
	HandlePasswordForgotten(to string) error
	NewEvent(ctx context.Context, users []*models.User, eventTime string) errsx.Map
	NewPayment(ctx context.Context, user *models.User, eventTime string) errsx.Map
	NewVote(ctx context.Context, user *models.User, eventTime string) errsx.Map
	PendingUser(ctx context.Context, user *models.User) errsx.Map
	SendRegistrationReminderEmail(ctx context.Context, user *models.User, registrationName string, daysLeft int) errsx.Map
	SendOTP(ctx context.Context, email, firstname string, otp *otpService.OTP) errsx.Map
	WelcomeUser(ctx context.Context, user *models.User) errsx.Map
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
		from:     "support@leviosa.care",
		email:    email,
		password: password,
	}, nil
}

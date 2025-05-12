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
	HandlePasswordForgotten(ctx context.Context, to string) error
	NewEvent(ctx context.Context, users []*models.User, eventTime string) error
	NewPayment(ctx context.Context, user *models.User, eventTime string) error
	NewVote(ctx context.Context, user *models.User, eventTime string) error
	SendRegistrationReminderEmail(ctx context.Context, user *models.User, registrationName string, daysLeft int) error
	SendOTP(ctx context.Context, email string, otp *otpService.OTP) error
	WelcomeUser(ctx context.Context, email string, user *models.User, legalAddress, companyInstagram string) error
}

type service struct {
	from     string
	email    string
	password string
	repo     Reader
}

func New(ctx context.Context, repo Reader) (Service, error) {
	email := os.Getenv("GMAIL_EMAIL")
	if email == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_EMAIL'"))
	}
	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_PASSWORD'"))
	}
	from, err := repo.GetCompanyEmail(ctx)
	if err != nil {
		return nil, domain.NewNotFoundErr(fmt.Errorf("environment variable 'GMAIL_PASSWORD'"))
	}
	return &service{
		// from:     "contact@leviosa.care",
		from:     from,
		email:    email,
		password: password,
	}, nil
}

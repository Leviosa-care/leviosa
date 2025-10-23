package ports

import "context"

// TODO: here I need to change the users used here to just the relevant information

type MailService interface {
	// main
	// HandlePasswordForgotten(ctx context.Context, to string) error
	// NewEvent(ctx context.Context, users []*models.User, eventTime string) error
	// NewPayment(ctx context.Context, user *models.User, eventTime string) error
	// NewVote(ctx context.Context, user *models.User, eventTime string) error
	// SendRegistrationReminderEmail(ctx context.Context, user *models.User, registrationName string, daysLeft int) error
	// SendOTP(ctx context.Context, email string, otp string) error
	WelcomeUser(ctx context.Context, email string) error
	// cache
	GetCompanyEmail(email string)
	SetCompanyEmail(email string)

	GetLogo(logo []byte)
	SetLogo(logo []byte)
}

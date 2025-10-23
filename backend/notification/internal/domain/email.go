package domain

import "time"

type EmailRequest struct {
	To           string
	Subject      string
	Template     string
	Data         interface{}
	CarbonCopy   map[string]string // email -> name
	Images       map[string]string // path -> rename
}

type OTPEmailData struct {
	OTP           string
	Year          int
	Address       string
	InstagramPath string
}

type WelcomeEmailData struct {
	Email         string
	Year          int
	Address       string
	InstagramPath string
}

type PasswordResetEmailData struct {
	Email         string
	Year          int
	Address       string
	InstagramPath string
}

type EventNotificationEmailData struct {
	EventTime     string
	Year          int
	Address       string
	InstagramPath string
}

type PaymentNotificationEmailData struct {
	EventTime     string
	Year          int
	Address       string
	InstagramPath string
}

type VoteNotificationEmailData struct {
	EventTime     string
	Year          int
	Address       string
	InstagramPath string
}

type RegistrationReminderEmailData struct {
	RegistrationName string
	DaysLeft         int
	Year             int
	Address          string
	InstagramPath    string
}

func NewOTPEmailData(otp, address, instagram string) *OTPEmailData {
	return &OTPEmailData{
		OTP:           otp,
		Year:          time.Now().Year(),
		Address:       address,
		InstagramPath: instagram,
	}
}

func NewWelcomeEmailData(email, address, instagram string) *WelcomeEmailData {
	return &WelcomeEmailData{
		Email:         email,
		Year:          time.Now().Year(),
		Address:       address,
		InstagramPath: instagram,
	}
}

func NewPasswordResetEmailData(email, address, instagram string) *PasswordResetEmailData {
	return &PasswordResetEmailData{
		Email:         email,
		Year:          time.Now().Year(),
		Address:       address,
		InstagramPath: instagram,
	}
}

func NewEventNotificationEmailData(eventTime, address, instagram string) *EventNotificationEmailData {
	return &EventNotificationEmailData{
		EventTime:     eventTime,
		Year:          time.Now().Year(),
		Address:       address,
		InstagramPath: instagram,
	}
}

func NewPaymentNotificationEmailData(eventTime, address, instagram string) *PaymentNotificationEmailData {
	return &PaymentNotificationEmailData{
		EventTime:     eventTime,
		Year:          time.Now().Year(),
		Address:       address,
		InstagramPath: instagram,
	}
}

func NewVoteNotificationEmailData(eventTime, address, instagram string) *VoteNotificationEmailData {
	return &VoteNotificationEmailData{
		EventTime:     eventTime,
		Year:          time.Now().Year(),
		Address:       address,
		InstagramPath: instagram,
	}
}

func NewRegistrationReminderEmailData(registrationName string, daysLeft int, address, instagram string) *RegistrationReminderEmailData {
	return &RegistrationReminderEmailData{
		RegistrationName: registrationName,
		DaysLeft:         daysLeft,
		Year:             time.Now().Year(),
		Address:          address,
		InstagramPath:    instagram,
	}
}
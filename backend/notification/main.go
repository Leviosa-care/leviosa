package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	"github.com/Leviosa-care/notification/internal/adapters/rabbitmq"
	"github.com/Leviosa-care/notification/internal/adapters/smtp"
	"github.com/Leviosa-care/notification/internal/adapters/twilio"
	"github.com/Leviosa-care/notification/internal/application"
	"github.com/Leviosa-care/notification/internal/application/email"
	"github.com/Leviosa-care/notification/internal/application/sms"
	"github.com/Leviosa-care/notification/internal/domain"
	"github.com/Leviosa-care/notification/internal/ports"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx := context.Background()

	service, err := initializeNotificationService(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize notification service: %v", err)
	}

	// Keep the service running
	select {}
}

func initializeNotificationService(ctx context.Context) (ports.NotificationService, error) {
	// Load environment variables
	gmailEmail := os.Getenv("GMAIL_EMAIL")
	if gmailEmail == "" {
		return nil, errs.NewNotFoundErr(fmt.Errorf("GMAIL_EMAIL environment variable is required"))
	}

	gmailPassword := os.Getenv("GMAIL_PASSWORD")
	if gmailPassword == "" {
		return nil, errs.NewNotFoundErr(fmt.Errorf("GMAIL_PASSWORD environment variable is required"))
	}

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if twilioAccountSID == "" {
		return nil, errs.NewNotFoundErr(fmt.Errorf("TWILIO_ACCOUNT_SID environment variable is required"))
	}

	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		return nil, errs.NewNotFoundErr(fmt.Errorf("TWILIO_AUTH_TOKEN environment variable is required"))
	}

	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	if twilioPhoneNumber == "" {
		return nil, errs.NewNotFoundErr(fmt.Errorf("TWILIO_PHONE_NUMBER environment variable is required"))
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://localhost:5672" // Default for development
	}

	// Initialize company cache with default values
	// These will be updated by the settings consumer
	cache := domain.NewCompanyCache("", "", "", []byte{})

	// Load initial settings from settings service
	if err := loadInitialSettings(ctx, cache); err != nil {
		log.Printf("Warning: Failed to load initial settings: %v", err)
	}

	// Initialize SMTP client
	smtpClient := smtp.New(gmailEmail, gmailPassword, "smtp.gmail.com", 587, cache)

	// Initialize Twilio client
	twilioClient, err := twilio.New(twilioAccountSID, twilioAuthToken, twilioPhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("initialize Twilio client: %w", err)
	}

	// Initialize application services
	emailService := email.New(smtpClient, cache)
	smsService := sms.New(twilioClient)

	// Initialize notification service
	notificationService := application.New(emailService, smsService)

	// Initialize RabbitMQ connection and consumer
	rabbitConn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, errs.NewConnectionFailureErr(fmt.Errorf("connect to RabbitMQ: %w", err))
	}

	ch, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, errs.NewConnectionFailureErr(fmt.Errorf("create RabbitMQ channel: %w", err))
	}

	settingsConsumer := rabbitmq.NewSettingsConsumer(cache)
	if err := settingsConsumer.StartSettingsConsumer(ctx, ch); err != nil {
		return nil, fmt.Errorf("start settings consumer: %w", err)
	}

	log.Println("Notification service initialized successfully")
	return notificationService, nil
}

func loadInitialSettings(ctx context.Context, cache *domain.CompanyCache) error {
	// TODO: Load initial settings from settings service HTTP API
	// This is a placeholder for now
	cache.SetCompanyEmail("noreply@leviosa.com")
	cache.SetCompanyLegalAddress("123 Rue Example, Paris, France")
	cache.SetCompanyInstagram("@leviosa_official")

	return nil
}


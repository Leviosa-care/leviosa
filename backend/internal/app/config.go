package app

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Environment string
	ServerPort  int

	// Database
	PostgresURL string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// RabbitMQ
	RabbitMQURL string

	// S3
	S3Endpoint        string
	S3Region          string
	S3BucketName      string
	S3AccessKeyID     string
	S3SecretAccessKey string

	// Vault
	VaultAddr  string
	VaultToken string

	// Encryption (encx)
	EncxKEKAlias    string
	EncxPepperAlias string

	// Stripe
	StripeSecretKey              string
	StripeWebhookSecret          string
	StripeConnectWebhookSecret   string

	// Twilio
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// OAuth
	GoogleClientID     string
	GoogleClientSecret string
	AppleClientID      string
	AppleClientSecret  string
	AppleTeamID        string
	AppleKeyID         string

	// Session
	SessionSecret string

	// Booking token signing secret (HMAC key for guest booking tokens)
	BookingTokenSecret string

	// Frontend origin for building public URLs (e.g. booking lookup links in SMS)
	FrontendOrigin string

	// Booking reminder scheduler
	ReminderIntervalMinutes  int
	ReminderWindowHours      int

	// Service Auth
	ServiceKeyCacheTTLSeconds int

	// Rate Limiting
	RateLimitSigninMax             int
	RateLimitSigninWindowSeconds   int
	RateLimitPasswordResetMax      int
	RateLimitPasswordResetWindowSeconds int
}

// LoadConfig loads configuration from environment variables
func LoadConfig(ctx context.Context) (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// Load .env file in development
	if env == "development" {
		if err := godotenv.Load("development.env"); err != nil {
			// Don't fail if file doesn't exist
			fmt.Printf("Warning: Could not load development.env: %v\n", err)
		}
	}

	cfg := &Config{
		Environment: env,
		ServerPort:  getEnvInt("SERVER_PORT", 8080),

		// Database
		PostgresURL: getEnv("DATABASE_URL", ""),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		// RabbitMQ
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),

		// S3
		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3Region:          getEnv("S3_REGION", "us-east-1"),
		S3BucketName:      getEnv("S3_BUCKET_NAME", ""),
		S3AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),

		// Vault
		VaultAddr:  getEnv("VAULT_ADDR", "http://localhost:8200"),
		VaultToken: getEnv("VAULT_TOKEN", ""),

		// Encryption (encx)
		EncxKEKAlias:    getEnv("ENCX_KEK_ALIAS", "leviosa-kek"),
		EncxPepperAlias: getEnv("ENCX_PEPPER_ALIAS", "leviosa"),

		// Stripe
		StripeSecretKey:            getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret:        getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripeConnectWebhookSecret: getEnv("STRIPE_CONNECT_WEBHOOK_SECRET", ""),

		// Twilio
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),

		// SMTP
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),

		// OAuth
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		AppleClientID:      getEnv("APPLE_CLIENT_ID", ""),
		AppleClientSecret:  getEnv("APPLE_CLIENT_SECRET", ""),
		AppleTeamID:        getEnv("APPLE_TEAM_ID", ""),
		AppleKeyID:         getEnv("APPLE_KEY_ID", ""),

		// Session
		SessionSecret: getEnv("SESSION_SECRET", "development-secret-key"),

		// Booking token
		BookingTokenSecret: getEnv("BOOKING_TOKEN_SECRET", ""),

		// Frontend
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),

		// Booking reminder scheduler
		ReminderIntervalMinutes: getEnvInt("REMINDER_INTERVAL_MINUTES", 15),
		ReminderWindowHours:    getEnvInt("REMINDER_WINDOW_HOURS", 24),

		// Service Auth
		ServiceKeyCacheTTLSeconds: getEnvInt("SERVICE_KEY_CACHE_TTL_SECONDS", 300),

		// Rate Limiting
		RateLimitSigninMax:                  getEnvInt("RATE_LIMIT_SIGNIN_MAX", 10),
		RateLimitSigninWindowSeconds:        getEnvInt("RATE_LIMIT_SIGNIN_WINDOW_SECONDS", 900),
		RateLimitPasswordResetMax:           getEnvInt("RATE_LIMIT_PASSWORD_RESET_MAX", 5),
		RateLimitPasswordResetWindowSeconds: getEnvInt("RATE_LIMIT_PASSWORD_RESET_WINDOW_SECONDS", 900),
	}

	// Validate required fields
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Validate session secret in non-development environments
	if cfg.Environment != "development" && cfg.SessionSecret == "development-secret-key" {
		return nil, fmt.Errorf("SESSION_SECRET must be set to a secure value in %s environment (not the development default)", cfg.Environment)
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

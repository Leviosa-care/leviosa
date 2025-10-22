package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/authuser/internal/ports"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	"github.com/hengadev/encx"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultOTPDuration    = 10 // Default OTP duration in minutes
	defaultOTPLength      = 6  // Default OTP length
	defaultOTPMaxAttempts = 3  // Default max attempts
)

type OTPService struct {
	repo   ports.OTPRepository
	crypto encx.CryptoService
	cache  ports.OTPCache
	mq     *amqp.Connection
}

func New(ctx context.Context, repo ports.OTPRepository, crypto encx.CryptoService, rabbitConn *amqp.Connection) (ports.OTPService, error) {
	// Initialize with default values first
	cache := domain.NewOTPCache(
		defaultOTPDuration,
		defaultOTPLength,
		defaultOTPMaxAttempts,
	)

	otpService := &OTPService{
		repo:   repo,
		crypto: crypto,
		cache:  cache,
		mq:     rabbitConn,
	}

	// Try to load current settings asynchronously
	go func() {
		if err := otpService.loadOTPSettings(ctx); err != nil {
			// TODO: Add proper logging here
			fmt.Printf("Failed to load OTP settings, using defaults: %v\n", err)
		}
	}()

	// Start RabbitMQ consumer for settings updates
	ch, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer channel: %w", err)
	}

	if err := otpService.StartOTPSettingConsumer(ctx, ch); err != nil {
		return nil, fmt.Errorf("failed to start settings consumer: %w", err)
	}

	return otpService, nil
}

func (s *OTPService) loadOTPSettings(ctx context.Context) error {
	keysList := []string{settings.OTPDuration, settings.OTPLength, settings.OTPMaxAttempts}
	baseURL := "http://backend:3500" // TODO: Make this configurable
	url := fmt.Sprintf("%s/settings/bulk?keys=%s", baseURL, strings.Join(keysList, ","))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch settings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("settings service returned status %d", resp.StatusCode)
	}

	var list []settings.SettingDTO
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return fmt.Errorf("failed to decode settings response: %w", err)
	}

	// Apply loaded settings
	for _, setting := range list {
		value, err := strconv.Atoi(setting.Value)
		if err != nil {
			fmt.Printf("Warning: Invalid setting value for %s: %s\n", setting.Key, setting.Value)
			continue
		}

		switch setting.Key {
		case settings.OTPDuration:
			s.SetOTPDuration(value)
		case settings.OTPLength:
			s.SetOTPLength(value)
		case settings.OTPMaxAttempts:
			s.SetOTPMaxAttempts(value)
		}
	}

	return nil
}

package otp

import (
	"context"
	"encoding/json"
	"fmt"

	mq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	sc "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Validation ranges for OTP settings
	minOTPDuration    = 1  // Minimum OTP duration in minutes
	maxOTPDuration    = 60 // Maximum OTP duration in minutes
	minOTPLength      = 4  // Minimum OTP code length in digits
	maxOTPLength      = 8  // Maximum OTP code length in digits
	minOTPMaxAttempts = 1  // Minimum max attempts
	maxOTPMaxAttempts = 10 // Maximum max attempts
)

// StartOTPSettingConsumer starts consuming OTP settings update messages from RabbitMQ.
//
// IMPORTANT: This consumer is currently NOT called during service initialization (see New() in service.go).
// OTP settings are hardcoded as constants for production reliability and simplicity.
//
// This infrastructure is preserved for potential future migration to a microservices architecture
// where dynamic runtime configuration updates may be needed.
//
// To re-enable this consumer:
// 1. Call this method from New() after service initialization
// 2. Ensure settings service is publishing updates to the OTP settings queue
// 3. Update application code to use s.GetOTPDuration(), s.GetOTPLength(), s.GetOTPMaxAttempts()
//
// The consumer validates all incoming setting values against defined ranges before applying them.
func (s *OTPService) StartOTPSettingConsumer(ctx context.Context, ch *amqp.Channel) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("retrieve logger in StartOTPSettingConsumer: %w", err)
	}
	msgs, err := ch.Consume(
		mq.OTPSettingsQueueName, // queue
		"",                      // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for d := range msgs {
			logger.InfoContext(ctx, "Received OTP settings update message",
				"body", string(d.Body))

			var payload rabbitmq.UpdatePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				logger.WarnContext(ctx, "Failed to unmarshal OTP settings message",
					"error", err,
					"body", string(d.Body))
				d.Nack(false, false) // Reject and don't requeue
				continue
			}

			switch payload.Key {
			case sc.OTPDuration:
				if duration := s.extractIntValue(payload.Value); duration != -1 {
					if err := s.validateOTPDuration(duration); err != nil {
						logger.WarnContext(ctx, "Invalid OTP duration value, ignoring update",
							"key", sc.OTPDuration,
							"value", duration,
							"error", err,
							"valid_range", fmt.Sprintf("%d-%d minutes", minOTPDuration, maxOTPDuration))
					} else {
						s.SetOTPDuration(duration)
						logger.InfoContext(ctx, "OTP duration updated",
							"key", sc.OTPDuration,
							"new_value", duration)
					}
				} else {
					logger.WarnContext(ctx, "Invalid type for OTP duration",
						"key", sc.OTPDuration,
						"type", fmt.Sprintf("%T", payload.Value),
						"value", payload.Value)
				}

			case sc.OTPLength:
				if length := s.extractIntValue(payload.Value); length != -1 {
					if err := s.validateOTPLength(length); err != nil {
						logger.WarnContext(ctx, "Invalid OTP length value, ignoring update",
							"key", sc.OTPLength,
							"value", length,
							"error", err,
							"valid_range", fmt.Sprintf("%d-%d digits", minOTPLength, maxOTPLength))
					} else {
						s.SetOTPLength(length)
						logger.InfoContext(ctx, "OTP length updated",
							"key", sc.OTPLength,
							"new_value", length)
					}
				} else {
					logger.WarnContext(ctx, "Invalid type for OTP length",
						"key", sc.OTPLength,
						"type", fmt.Sprintf("%T", payload.Value),
						"value", payload.Value)
				}

			case sc.OTPMaxAttempts:
				if maxAttempts := s.extractIntValue(payload.Value); maxAttempts != -1 {
					if err := s.validateOTPMaxAttempts(maxAttempts); err != nil {
						logger.WarnContext(ctx, "Invalid OTP max attempts value, ignoring update",
							"key", sc.OTPMaxAttempts,
							"value", maxAttempts,
							"error", err,
							"valid_range", fmt.Sprintf("%d-%d attempts", minOTPMaxAttempts, maxOTPMaxAttempts))
					} else {
						s.SetOTPMaxAttempts(maxAttempts)
						logger.InfoContext(ctx, "OTP max attempts updated",
							"key", sc.OTPMaxAttempts,
							"new_value", maxAttempts)
					}
				} else {
					logger.WarnContext(ctx, "Invalid type for OTP max attempts",
						"key", sc.OTPMaxAttempts,
						"type", fmt.Sprintf("%T", payload.Value),
						"value", payload.Value)
				}

			default:
				logger.WarnContext(ctx, "Received unknown OTP settings key",
					"key", payload.Key,
					"value", payload.Value)
			}

			d.Ack(false) // Acknowledge the message after processing
		}
	}()

	logger.InfoContext(ctx, "OTP settings consumer started successfully")
	return nil
}

// extractIntValue safely extracts integer values from interface{} that might come as float64 from JSON
func (s *OTPService) extractIntValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// Try to parse string to int if needed
		if parsed, err := fmt.Sscanf(v, "%d", new(int)); err == nil && parsed == 1 {
			var result int
			fmt.Sscanf(v, "%d", &result)
			return result
		}
	}
	return -1 // Invalid value
}

// validateOTPDuration validates that OTP duration is within acceptable range
func (s *OTPService) validateOTPDuration(duration int) error {
	if duration < minOTPDuration || duration > maxOTPDuration {
		return fmt.Errorf("duration must be between %d and %d minutes, got %d", minOTPDuration, maxOTPDuration, duration)
	}
	return nil
}

// validateOTPLength validates that OTP length is within acceptable range
func (s *OTPService) validateOTPLength(length int) error {
	if length < minOTPLength || length > maxOTPLength {
		return fmt.Errorf("length must be between %d and %d digits, got %d", minOTPLength, maxOTPLength, length)
	}
	return nil
}

// validateOTPMaxAttempts validates that OTP max attempts is within acceptable range
func (s *OTPService) validateOTPMaxAttempts(maxAttempts int) error {
	if maxAttempts < minOTPMaxAttempts || maxAttempts > maxOTPMaxAttempts {
		return fmt.Errorf("max attempts must be between %d and %d, got %d", minOTPMaxAttempts, maxOTPMaxAttempts, maxAttempts)
	}
	return nil
}

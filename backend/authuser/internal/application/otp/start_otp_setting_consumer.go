package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	sc "github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// StartConsumer starts consuming messages from the mail settings queue.
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
			log.Printf("Received message: %s", d.Body)
			var payload rabbitmq.UpdatePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				logger.WarnContext(ctx, fmt.Sprintf("failed to unmarshal message: %s", err))
				d.Nack(false, false) // Reject and don't requeue (for now, handle errors)
				continue
			}

			switch payload.Key {
			case sc.OTPDuration:
				if duration := s.extractIntValue(payload.Value); duration != -1 {
					s.SetOTPDuration(duration)
					logger.InfoContext(ctx, fmt.Sprintf("OTP service cache updated: %s = %d", sc.OTPDuration, duration))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T (value: %v)", sc.OTPDuration, payload.Value, payload.Value))
				}
			case sc.OTPLength:
				if length := s.extractIntValue(payload.Value); length != -1 {
					s.SetOTPLength(length)
					logger.InfoContext(ctx, fmt.Sprintf("OTP service cache updated: %s = %d", sc.OTPLength, length))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T (value: %v)", sc.OTPLength, payload.Value, payload.Value))
				}
			case sc.OTPMaxAttempts:
				if maxAttempts := s.extractIntValue(payload.Value); maxAttempts != -1 {
					s.SetOTPMaxAttempts(maxAttempts)
					logger.InfoContext(ctx, fmt.Sprintf("OTP service cache updated: %s = %d", sc.OTPMaxAttempts, maxAttempts))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T (value: %v)", sc.OTPMaxAttempts, payload.Value, payload.Value))
				}
			default:
				log.Printf("received unknown settings update: %v", payload)
			}
			d.Ack(false) // Acknowledge the message after processing
		}
	}()

	logger.InfoContext(ctx, "OTPSettings service consumer started.")
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

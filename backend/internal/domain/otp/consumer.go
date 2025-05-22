package otpService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hengadev/leviosa/internal/broker/rabbitmq"
	"github.com/hengadev/leviosa/internal/domain/settings"
	"github.com/hengadev/leviosa/pkg/ctxutil"

	amqp "github.com/rabbitmq/amqp091-go"
)

// StartConsumer starts consuming messages from the mail settings queue.
func (s *service) StartOTPSettingConsumer(ctx context.Context, ch *amqp.Channel) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("retrieve logger in StartOTPSettingConsumer: %w", err)
	}
	msgs, err := ch.Consume(
		rabbitmq.OTPSettingsQueueName, // queue
		"",                            // consumer
		false,                         // auto-ack
		false,                         // exclusive
		false,                         // no-local
		false,                         // no-wait
		nil,                           // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			var payload rabbitmq.SettingsUpdatePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				logger.WarnContext(ctx, fmt.Sprintf("failed to unmarshal message: %v", err))
				d.Nack(false, false) // Reject and don't requeue (for now, handle errors)
				continue
			}

			switch payload.Key {
			case settings.OTPDurationKey:
				if duration, ok := payload.Value.(int); ok {
					s.SetOTPDuration(duration)
					logger.InfoContext(ctx, fmt.Sprintf("OTP service cache updated: %s = %s", settings.OTPDurationKey, err))
				} else {
					logger.InfoContext(ctx, fmt.Sprintf("invalid type for %s: %T", settings.OTPDurationKey, payload.Value))
				}
			case settings.OTPLengthKey:
				if length, ok := payload.Value.(int); ok {
					s.SetOTPLength(length)
					logger.InfoContext(ctx, fmt.Sprintf("OTP service cache updated: %s = %s", settings.OTPLengthKey, err))
				} else {
					logger.InfoContext(ctx, fmt.Sprintf("invalid type for %s: %T", settings.OTPLengthKey, payload.Value))
				}
			case settings.OTPMaxAttemptsKey:
				if maxAttempts, ok := payload.Value.(int); ok {
					s.SetOTPMaxAttempts(maxAttempts)
					logger.InfoContext(ctx, fmt.Sprintf("OTP service cache updated: %s = %s", settings.OTPMaxAttemptsKey, err))
				} else {
					logger.InfoContext(ctx, fmt.Sprintf("invalid type for %s: %T", settings.OTPMaxAttemptsKey, payload.Value))
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

package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	rbmq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	sc "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO: change that function to make it work properly with the right value from contact etc..

// StartSMSSettingConsumer starts consuming messages from the mail settings queue.
func (s *SMSService) StartSMSSettingConsumer(ctx context.Context, ch *amqp.Channel) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("retrieve logger in StartMailSettingConsumer: %w", err)
	}
	msgs, err := ch.Consume(
		// rabbitmq.MailSettingsQueueName, // queue
		rbmq.NotificationSettingsQueueName, // queue
		"",                                 // consumer
		false,                              // auto-ack
		false,                              // exclusive
		false,                              // no-local
		false,                              // no-wait
		nil,                                // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			var payload rabbitmq.SettingsUpdatePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				logger.WarnContext(ctx, fmt.Sprintf("failed to unmarshal message: %s", err))
				d.Nack(false, false) // Reject and don't requeue (for now, handle errors)
				continue
			}

			switch payload.Key {
			case sc.CompanyPhone:
				if phone, ok := payload.Value.(string); ok {
					s.SetPhone(phone)
					logger.InfoContext(ctx, fmt.Sprintf("SMS service cache updated: %s = %s", sc.CompanyPhone, err))
				} else {
					logger.InfoContext(ctx, fmt.Sprintf("invalid type for %s: %T", sc.CompanyPhone, payload.Value))
				}
			default:
				log.Printf("received unknown settings update: %v", payload)
			}
			d.Ack(false) // Acknowledge the message after processing
		}
	}()

	logger.InfoContext(ctx, "SMSSettings service consumer started.")

	return nil
}

package mailService

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
func (s *service) StartMailSettingConsumer(ctx context.Context, ch *amqp.Channel) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("retrieve logger in StartMailSettingConsumer: %w", err)
	}
	msgs, err := ch.Consume(
		rabbitmq.MailSettingsQueueName, // queue
		"",                             // consumer
		false,                          // auto-ack
		false,                          // exclusive
		false,                          // no-local
		false,                          // no-wait
		nil,                            // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			var payload rabbitmq.SettingsUpdatePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				logger.WarnContext(ctx, "failed to unmarshal message: %v", err)
				d.Nack(false, false) // Reject and don't requeue (for now, handle errors)
				continue
			}

			switch payload.Key {
			case settings.CompanyEmailKey:
				if email, ok := payload.Value.(string); ok {
					s.SetCompanyEmail(email)
					logger.InfoContext(ctx, fmt.Sprintf("Mail service cache updated: %s = %s", settings.CompanyEmailKey, err))
				} else {
					logger.InfoContext(ctx, fmt.Sprintf("invalid type for %s: %T", settings.CompanyEmailKey, payload.Value))
				}
			case settings.CompanyLogoKey:
				if logo, ok := payload.Value.([]byte); ok {
					s.SetLogo(logo)
					logger.InfoContext(ctx, fmt.Sprintf("Mail service cache updated: %s = %s", settings.CompanyLogoKey, logo))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T", settings.CompanyLogoKey, payload.Value))
				}
			default:
				log.Printf("received unknown settings update: %v", payload)
			}
			d.Ack(false) // Acknowledge the message after processing
		}
	}()

	logger.InfoContext(ctx, "MailSettings service consumer started.")
	return nil
}

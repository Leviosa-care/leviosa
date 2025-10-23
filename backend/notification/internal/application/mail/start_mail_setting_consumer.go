package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	sc "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// StartConsumer starts consuming messages from the mail settings queue.
func (s *MailService) StartMailSettingConsumer(ctx context.Context, ch *amqp.Channel) error {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("retrieve logger in StartMailSettingConsumer: %w", err)
	}
	msgs, err := ch.Consume(
		// rabbitmq.MailSettingsQueueName, // queue
		mq.NotificationSettingsQueueName, // queue
		"",                               // consumer
		false,                            // auto-ack
		false,                            // exclusive
		false,                            // no-local
		false,                            // no-wait
		nil,                              // args
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
			case sc.CompanyEmail:
				if email, ok := payload.Value.(string); ok {
					s.SetCompanyEmail(email)
					logger.InfoContext(ctx, fmt.Sprintf("Mail service cache updated: %s = %s", sc.CompanyEmail, err))
				} else {
					logger.InfoContext(ctx, fmt.Sprintf("invalid type for %s: %T", sc.CompanyEmail, payload.Value))
				}
			case sc.CompanyLogo:
				if logo, ok := payload.Value.([]byte); ok {
					s.SetLogo(logo)
					logger.InfoContext(ctx, fmt.Sprintf("Mail service cache updated: %s = %s", sc.CompanyLogo, logo))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T", sc.CompanyLogo, payload.Value))
				}
			case sc.CompanyInstagram:
				if insta, ok := payload.Value.(string); ok {
					s.SetCompanyInstagram(insta)
					logger.InfoContext(ctx, fmt.Sprintf("Mail service cache updated: %s = %s", sc.CompanyInstagram, insta))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T", sc.CompanyInstagram, payload.Value))
				}
			case sc.CompanyLegalAddress:
				if addr, ok := payload.Value.(string); ok {
					s.SetCompanyLegalAddress(addr)
					logger.InfoContext(ctx, fmt.Sprintf("Mail service cache updated: %s = %s", sc.CompanyLegalAddress, addr))
				} else {
					logger.WarnContext(ctx, fmt.Sprintf("invalid type for %s: %T", sc.CompanyLegalAddress, payload.Value))
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

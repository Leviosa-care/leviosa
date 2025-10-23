package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	"github.com/Leviosa-care/notification/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SettingsConsumer struct {
	cache *domain.CompanyCache
}

func NewSettingsConsumer(cache *domain.CompanyCache) *SettingsConsumer {
	return &SettingsConsumer{
		cache: cache,
	}
}

func (c *SettingsConsumer) StartSettingsConsumer(ctx context.Context, ch *amqp.Channel) error {
	if ch == nil {
		return errs.NewInvalidValueErr(("RabbitMQ channel cannot be nil"))
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
		return errs.NewConnectionFailureErr(fmt.Errorf("start RabbitMQ settings consumer: %w", err))
	}

	go func() {
		for d := range msgs {
			if err := c.processMessage(ctx, d.Body); err != nil {
				d.Nack(false, false)
				continue
			}
			d.Ack(false)
		}
	}()

	return nil
}

func (c *SettingsConsumer) processMessage(ctx context.Context, body []byte) error {
	var payload rabbitmq.SettingsUpdatePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return errs.NewInvalidValueErr(fmt.Errorf("unmarshal settings update payload: %w", err))
	}

	switch payload.Key {
	case settings.CompanyEmail:
		if email, ok := payload.Value.(string); ok {
			c.cache.SetCompanyEmail(email)
		} else {
			return errs.NewInvalidValueErr(fmt.Errorf("invalid type for %s: expected string, got %T", settings.CompanyEmail, payload.Value))
		}

	case settings.CompanyLegalAddress:
		if address, ok := payload.Value.(string); ok {
			c.cache.SetCompanyLegalAddress(address)
		} else {
			return errs.NewInvalidValueErr(fmt.Errorf("invalid type for %s: expected string, got %T", settings.CompanyLegalAddress, payload.Value))
		}

	case settings.CompanyInstagram:
		if instagram, ok := payload.Value.(string); ok {
			c.cache.SetCompanyInstagram(instagram)
		} else {
			return errs.NewInvalidValueErr(fmt.Errorf("invalid type for %s: expected string, got %T", settings.CompanyInstagram, payload.Value))
		}

	default:
	}

	return nil
}

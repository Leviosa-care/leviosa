package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implements the EventPublisher interface using RabbitMQ.
type RabbitMQPublisher struct {
	conn *amqp.Connection
}

// NewPublisher creates a new RabbitMQ event publisher.
func NewPublisher(conn *amqp.Connection) ports.EventPublisher {
	return &RabbitMQPublisher{
		conn: conn,
	}
}

// PublishSettingUpdate publishes a setting update event to RabbitMQ.
func (p *RabbitMQPublisher) PublishSettingUpdate(ctx context.Context, key string, value any) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	payload := rabbitmq.UpdatePayload{
		Key:   key,
		Value: value,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	err = ch.PublishWithContext(
		ctx,
		mq.SettingsExchangeName, // exchange
		mq.SettingsRoutingKey,   // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published settings update: key=%s, value=%v", key, value)
	return nil
}

package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// TokenDurationCache interface for updating cached token durations
type TokenDurationCache interface {
	UpdateAccessTokenDuration(ctx context.Context, durationMinutes int) error
	UpdateRefreshTokenDuration(ctx context.Context, durationHours int) error
}

// SettingsConsumer handles settings update messages
type SettingsConsumer struct {
	conn  *amqp.Connection
	cache TokenDurationCache
}

// NewSettingsConsumer creates a new settings consumer
func NewSettingsConsumer(conn *amqp.Connection, cache TokenDurationCache) *SettingsConsumer {
	return &SettingsConsumer{
		conn:  conn,
		cache: cache,
	}
}

// Start begins consuming settings update messages
func (c *SettingsConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare settings exchange (should already exist from settings service)
	if err := rabbitmq.DeclareExchange(ch, mq.SettingsExchangeName, "direct"); err != nil {
		return err
	}

	// Declare authuser settings queue
	if err := rabbitmq.DeclareQueue(ch, mq.AuthUserSettingsQueueName); err != nil {
		return err
	}

	// Bind queue to exchange
	if err := rabbitmq.BindQueue(ch, mq.AuthUserSettingsQueueName, mq.SettingsRoutingKey, mq.SettingsExchangeName); err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		mq.AuthUserSettingsQueueName, // queue
		"authuser-settings-consumer", // consumer
		false,                        // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		return err
	}

	log.Printf("Starting authuser settings consumer...")

	// Process messages
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping authuser settings consumer...")
			return ctx.Err()
		case msg := <-msgs:
			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("Error processing settings message: %v", err)
				// Nack the message so it can be retried
				msg.Nack(false, true)
			} else {
				// Ack the message
				msg.Ack(false)
			}
		}
	}
}

// processMessage handles individual settings update messages
func (c *SettingsConsumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	var payload rabbitmq.UpdatePayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		return err
	}

	log.Printf("Received settings update: key=%s, value=%v", payload.Key, payload.Value)

	switch payload.Key {
	case settings.AccessTokenDuration:
		if duration, ok := payload.Value.(float64); ok {
			return c.cache.UpdateAccessTokenDuration(ctx, int(duration))
		}
	case settings.RefreshTokenDuration:
		if duration, ok := payload.Value.(float64); ok {
			return c.cache.UpdateRefreshTokenDuration(ctx, int(duration))
		}
	default:
		// Ignore other settings updates
		log.Printf("Ignoring settings update for key: %s", payload.Key)
	}

	return nil
}

// SetupSettingsConsumer sets up the RabbitMQ infrastructure for settings consumption
func SetupSettingsConsumer(ctx context.Context, ch *amqp.Channel) error {
	// Declare settings exchange
	if err := rabbitmq.DeclareExchange(ch, mq.SettingsExchangeName, "direct"); err != nil {
		return err
	}

	// Declare authuser settings queue
	if err := rabbitmq.DeclareQueue(ch, mq.AuthUserSettingsQueueName); err != nil {
		return err
	}

	// Bind queue to exchange
	if err := rabbitmq.BindQueue(ch, mq.AuthUserSettingsQueueName, mq.SettingsRoutingKey, mq.SettingsExchangeName); err != nil {
		return err
	}

	return nil
}

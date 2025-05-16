package rabbitmq

import (
	"context"
	"fmt"

	"github.com/hengadev/errsx"
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupRabbitMQ declares the necessary exchange and queue for settings updates.
func SetupRabbitMQ(ctx context.Context, ch *amqp.Channel) error {
	var errs errsx.Map
	if err := SetupSettingsServiceQueues(ctx, ch); err != nil {
		errs.Set("setup settings service queues", err)
	}
	return errs.AsError()
}

// SetupSettingsServiceQueues declares the exchange and queues necessary for
// services (like mail and OTP) to receive settings updates.
func SetupSettingsServiceQueues(ctx context.Context, ch *amqp.Channel) error {
	// Declare the exchange for settings updates
	err := declareExchange(ch, SettingsExchangeName, "direct")
	if err != nil {
		return err
	}

	// Declare the queue for mail service settings updates
	err = declareQueue(ch, MailSettingsQueueName)
	if err != nil {
		return err
	}

	// Bind the mail settings queue to the exchange
	err = bindQueue(ch, MailSettingsQueueName, SettingsRoutingKey, SettingsExchangeName)
	if err != nil {
		return err
	}

	// Declare the queue for OTP service settings updates
	err = declareQueue(ch, OTPSettingsQueueName)
	if err != nil {
		return err
	}

	// Bind the OTP settings queue to the exchange (using the same routing key for now)
	err = bindQueue(ch, OTPSettingsQueueName, SettingsRoutingKey, SettingsExchangeName)
	if err != nil {
		return err
	}

	return nil
}

// declareExchange is a helper function to declare an exchange.
func declareExchange(ch *amqp.Channel, name, kind string) error {
	err := ch.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange '%s': %w", name, err)
	}
	return nil
}

// declareQueue is a helper function to declare a queue.
func declareQueue(ch *amqp.Channel, name string) error {
	_, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue '%s': %w", name, err)
	}
	return nil
}

// bindQueue is a helper function to bind a queue to an exchange.
func bindQueue(ch *amqp.Channel, queueName, routingKey, exchangeName string) error {
	err := ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue '%s' to exchange '%s' with routing key '%s': %w", queueName, exchangeName, routingKey, err)
	}
	return nil
}

// NewChannel creates a new RabbitMQ channel.
func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	return ch, nil
}

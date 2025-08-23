package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeclareQueue is a helper function to declare a queue.
func DeclareQueue(ch *amqp.Channel, name string) error {
	// _, err := ch.QueueDeclare(
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

// BindQueue is a helper function to bind a queue to an exchange.
func BindQueue(ch *amqp.Channel, queueName, routingKey, exchangeName string) error {
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

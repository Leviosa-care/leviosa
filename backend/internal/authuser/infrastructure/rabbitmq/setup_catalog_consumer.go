package rabbitmq

import (
	"context"
	"fmt"

	rabbitmqContracts "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupCatalogConsumer sets up the RabbitMQ infrastructure for catalog event consumption
func SetupCatalogConsumer(ctx context.Context, ch *amqp.Channel) error {
	// Declare catalog exchange
	if err := rabbitmq.DeclareExchange(ch, rabbitmqContracts.CatalogExchangeName, "topic"); err != nil {
		return fmt.Errorf("declare catalog exchange: %w", err)
	}

	// Declare authuser catalog queue
	if err := rabbitmq.DeclareQueue(ch, rabbitmqContracts.AuthUserCatalogQueueName); err != nil {
		return fmt.Errorf("declare authuser catalog queue: %w", err)
	}

	// Bind queue to exchange for all catalog events
	routingKeys := []string{
		rabbitmqContracts.CategoryCreatedRoutingKey,
		rabbitmqContracts.CategoryUpdatedRoutingKey,
		rabbitmqContracts.CategoryDeletedRoutingKey,
		rabbitmqContracts.ProductCreatedRoutingKey,
		rabbitmqContracts.ProductUpdatedRoutingKey,
		rabbitmqContracts.ProductDeletedRoutingKey,
	}

	for _, routingKey := range routingKeys {
		if err := rabbitmq.BindQueue(ch, rabbitmqContracts.AuthUserCatalogQueueName, routingKey, rabbitmqContracts.CatalogExchangeName); err != nil {
			return fmt.Errorf("bind queue to %s: %w", routingKey, err)
		}
	}

	return nil
}
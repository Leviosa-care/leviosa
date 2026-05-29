package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Setup(ctx context.Context, ch *amqp.Channel) error {
	// Setup catalog consumer queues
	if err := SetupCatalogConsumer(ctx, ch); err != nil {
		return err
	}

	return nil
}

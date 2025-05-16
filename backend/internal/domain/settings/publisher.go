package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hengadev/leviosa/internal/broker/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher interface defines the method for publishing settings updates.
type Publisher interface {
	PublishSettingUpdate(ctx context.Context, key string, value any) error
}

// service (assuming your settings service struct is named 'service') now implements the Publisher interface.
func (s *service) PublishSettingUpdate(ctx context.Context, key string, value any) error {
	ch, err := s.rabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	payload := rabbitmq.SettingsUpdatePayload{
		Key:   key,
		Value: value,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	err = ch.PublishWithContext(
		ctx,
		rabbitmq.SettingsExchangeName, // exchange
		rabbitmq.SettingsRoutingKey,   // routing key
		false,                         // mandatory
		false,                         // immediate
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

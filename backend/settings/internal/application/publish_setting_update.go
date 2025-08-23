package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (s *SettingsService) PublishSettingUpdate(ctx context.Context, key string, value any) error {
	ch, err := s.mq.Channel()
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

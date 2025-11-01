package rabbitmq

import (
	"context"

	mq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Setup(ctx context.Context, ch *amqp.Channel) error {
	if err := rabbitmq.DeclareExchange(ch, mq.SettingsExchangeName, "direct"); err != nil {
		return err
	}
	// Declare the queue for mail service settings updates
	if err := rabbitmq.DeclareQueue(ch, mq.NotificationSettingsQueueName); err != nil {
		return err
	}

	// Bind the mail settings queue to the exchange
	if err := rabbitmq.BindQueue(ch, mq.NotificationSettingsQueueName, mq.SettingsRoutingKey, mq.SettingsExchangeName); err != nil {
		return err
	}

	// Declare the queue for OTP service settings updates
	if err := rabbitmq.DeclareQueue(ch, mq.OTPSettingsQueueName); err != nil {
		return err
	}

	// Bind the OTP settings queue to the exchange (using the same routing key for now)
	if err := rabbitmq.BindQueue(ch, mq.OTPSettingsQueueName, mq.SettingsRoutingKey, mq.SettingsExchangeName); err != nil {
		return err
	}
	return nil
}

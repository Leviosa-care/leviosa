package rabbitmq

import (
	"context"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Setup(ctx context.Context, ch *amqp.Channel) error {
	// OTP notification exchange and queues
	if err := rabbitmq.DeclareExchange(ch, mq.OTPNotificationExchangeName, "direct"); err != nil {
		return err
	}

	// Declare the queue for OTP email notifications
	if err := rabbitmq.DeclareQueue(ch, mq.OTPEmailQueueName); err != nil {
		return err
	}

	// Bind the OTP email queue to the exchange
	if err := rabbitmq.BindQueue(ch, mq.OTPEmailQueueName, mq.OTPEmailRoutingKey, mq.OTPNotificationExchangeName); err != nil {
		return err
	}

	// Declare the queue for OTP SMS notifications
	if err := rabbitmq.DeclareQueue(ch, mq.OTPSMSQueueName); err != nil {
		return err
	}

	// Bind the OTP SMS queue to the exchange
	if err := rabbitmq.BindQueue(ch, mq.OTPSMSQueueName, mq.OTPSMSRoutingKey, mq.OTPNotificationExchangeName); err != nil {
		return err
	}

	return nil
}

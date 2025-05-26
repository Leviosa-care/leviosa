package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hengadev/leviosa/internal/broker/rabbitmq"
	"github.com/hengadev/leviosa/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func setBroker(ctx context.Context, conf *config.RabbitSecrets) (*amqp.Connection, error) {
	amqpURL := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(conf.User, conf.Password),
		Host:   fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Path:   "/",
	}
	rabbitConn, err := amqp.Dial(amqpURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	setupChan, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}
	defer setupChan.Close()

	if err := rabbitmq.SetupRabbitMQ(ctx, setupChan); err != nil {
		rabbitConn.Close()
		return nil, fmt.Errorf("failed to setup RabbitMQ: %w", err)
	}
	return rabbitConn, nil
}

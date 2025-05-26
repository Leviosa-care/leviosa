package config

import (
	cfg "github.com/hengadev/leviosa/pkg/config"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/hengadev/errsx"
)

func (c *Config) GetRabbitMQ() *cfg.RabbitSecrets {
	return c.rabbitmq
}

// RabbitMQDefault returns a default configuration for RabbitMQ.
func RabbitMQDefault() *cfg.RabbitSecrets {
	return &cfg.RabbitSecrets{
		Host:     "localhost", // Change if your RabbitMQ is running elsewhere
		Port:     "5672",      // Default RabbitMQ port
		User:     "guest",     // Default user (change in production)
		Password: "guest",     // Default password (change in production)
	}
}

func (c *Config) setRabbitMQ(env envmode.Mode) error {
	var host, port, user, password string
	var errs errsx.Map
	switch env {
	case envmode.Dev:
		c.rabbitmq = RabbitMQDefault()
	case envmode.Prod, envmode.Staging:

		host = c.viper.GetString("rabbitmq.host")
		port = c.viper.GetString("rabbitmq.port")
		user = c.viper.GetString("rabbitmq.user")
		password = c.viper.GetString("rabbitmq.password")
	}
	if host == "" {
		errs.Set("RABBITMQ_HOST", "'RABBITMQ_HOST' environment variable not set; please define it to specify RabbitMQ host")
	}
	if port == "" {
		errs.Set("RABBITMQ_PORT", "'RABBITMQ_PORT' environment variable not set; please define it to specify RabbitMQ port")
	}
	if user == "" {
		errs.Set("RABBITMQ_USER", "'RABBITMQ_USER' environment variable not set; please define it to specify RabbitMQ user")
	}
	if password == "" {
		errs.Set("RABBITMQ_PASSWORD", "'RABBITMQ_PASSWORD' environment variable not set; please define it to specify RabbitMQ password")
	}
	c.rabbitmq.Host = host
	c.rabbitmq.Port = port
	c.rabbitmq.User = user
	c.rabbitmq.Password = password
	return errs.AsError()
}

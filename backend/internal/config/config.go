package config

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/hengadev/leviosa/pkg/config"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/hengadev/errsx"
	"github.com/spf13/viper"
)

type Config struct {
	viper    *viper.Viper
	postgres *cfg.PostgresSecrets
	redis    *cfg.RedisSecrets
	s3       *cfg.S3Secrets
	rabbitmq *cfg.RabbitSecrets
}

func Load(ctx context.Context, mode envmode.Mode) (*Config, error) {
	v := viper.New()
	if mode == envmode.Dev {
		v.AddConfigPath(".")
		v.SetConfigName(mode.String())
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("viper reading :", err)
		}
	}
	c := &Config{
		viper: v,
	}
	envVarsToKeys := map[string]struct {
		required bool
		key      string
	}{
		"REDIS_ADDR":     {required: true, key: "redis.addr"},
		"REDIS_DB":       {required: true, key: "redis.db"},
		"REDIS_PASSWORD": {required: true, key: "redis.password"},

		"STRIPE_SECRET_KEY": {required: true, key: "stripe.secret.key"},

		"GMAIL_EMAIL":    {required: true, key: "gmail.email"},
		"GMAIL_PASSWORD": {required: true, key: "gmail.password"},

		"AWS_REGION":            {required: true, key: "aws.region"},
		"AWS_ACCESS_KEY_ID":     {required: true, key: "aws.access.key.id"},
		"AWS_SECRET_ACCESS_KEY": {required: true, key: "aws.secret.access.key"},

		"LOGGING_SALT": {required: true, key: "logging.salt"},

		"BUCKETNAME": {required: true, key: "s3.bucketname"},

		"RABBITMQ_HOST":     {required: true, key: "rabbitmq.host"},
		"RABBITMQ_PORT":     {required: true, key: "rabbitmq.port"},
		"RABBITMQ_USER":     {required: true, key: "rabbitmq.user"},
		"RABBITMQ_PASSWORD": {required: true, key: "rabbitmq.password"},
	}
	var errs errsx.Map
	for envVar, requiredKey := range envVarsToKeys {
		if os.Getenv(envVar) == "" && requiredKey.required == true {
			errs.Set("get environment variable", fmt.Errorf("missing required env variables: %s", envVar))
		}
		if err := c.viper.BindEnv(requiredKey.key, envVar); err != nil {
			errs.Set("bind environment variable", fmt.Errorf("bind env: %w", err))
		}
	}
	if err := c.setPostgres(mode); err != nil {
		errs.Set("postgres configuration", fmt.Errorf("set PostgreSQL: %w", err))
	}
	if err := c.setRedis(mode); err != nil {
		errs.Set("redis configuration", fmt.Errorf("set Redis: %w", err))
	}
	if err := c.setS3(mode); err != nil {
		errs.Set("S3 configuration", fmt.Errorf("set S3: %w", err))
	}
	if err := c.setRabbitMQ(mode); err != nil {
		errs.Set("Rabbit MQ configuration", fmt.Errorf("set S3: %w", err))
	}
	return c, errs.AsError()
}

package config

import (
	"fmt"

	cfg "github.com/hengadev/leviosa/pkg/config"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/hengadev/errsx"
)

func (c *Config) GetPostgres() *cfg.PostgresSecrets {
	return c.postgres
}

func (c *Config) setPostgres(env envmode.Mode) error {
	var errs errsx.Map
	var host string
	switch env {
	case envmode.Dev:
		host = "localhost"
	case envmode.Staging, envmode.Prod:
		host = c.viper.GetString("postgres.host")
	}
	c.postgres.Host = host
	c.postgres.User = c.viper.GetString("postgres.user")
	c.postgres.Password = c.viper.GetString("postgres.password")
	c.postgres.Port = c.viper.GetInt("postgres.port")
	c.postgres.DB = fmt.Sprintf("%s_%s", env.String(), c.viper.GetString("postgres.db"))
	return errs.AsError()
}

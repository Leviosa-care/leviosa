package config

import (
	cfg "github.com/hengadev/leviosa/pkg/config"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/hengadev/errsx"
)

type postgresCreds struct {
	name     string
	user     string
	password string
	port     int
}

func (c *Config) GetPostgres() *cfg.PostgresSecrets {
	return c.postgres
}

func (c *Config) setPostgres(env envmode.Mode) error {
	var errs errsx.Map
	return errs.AsError()
}

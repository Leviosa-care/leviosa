package config

import (
	"github.com/hengadev/errsx"
)

type PostgresSecrets struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func (p PostgresSecrets) Validate() error {
	var errs errsx.Map
	if p.Host == "" {
		errs.Set("postgres host", "postgres host cannot be empty")
	}
	if p.User == "" {
		errs.Set("postgres user", "postgres user cannot be empty")
	}
	if p.Port <= 0 || p.Port > 65535 {
		errs.Set("postgres port", "postgres port must be between 1 and 65535")
	}
	if p.Password == "" {
		errs.Set("postgres password", "postgres password must be non-negative")
	}
	return errs.AsError()
}

func (p PostgresSecrets) GetType() string {
	return "postgres"
}

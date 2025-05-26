package config

import (
	"github.com/hengadev/errsx"
)

type PostgresSecrets struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func (p PostgresSecrets) Validate() error {
	var errs errsx.Map
	if p.Name == "" {
		errs.Set("postgres name", "postgres name cannot be empty")
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

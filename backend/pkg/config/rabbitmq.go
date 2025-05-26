package config

import (
	"strconv"

	"github.com/hengadev/errsx"
)

type RabbitSecrets struct {
	Host     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func (r RabbitSecrets) Validate() error {
	var errs errsx.Map
	if r.Host == "" {
		errs.Set("rabbit name", "rabbit name cannot be empty")
	}
	if r.User == "" {
		errs.Set("rabbit user", "rabbit user cannot be empty")
	}
	port, err := strconv.Atoi(r.Port)
	if err != nil {
		errs.Set("rabbit port", "rabbit port cannot be convert to int")
	}
	if port <= 0 || port > 65535 {
		errs.Set("rabbit port", "rabbit port must be between 1 and 65535")
	}
	if r.Password == "" {
		errs.Set("rabbit password", "rabbit password must be non-negative")
	}
	return errs.AsError()
}

func (r RabbitSecrets) GetType() string {
	return "rabbit"
}

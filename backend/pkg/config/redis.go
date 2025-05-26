package config

import (
	"github.com/hengadev/errsx"
)

type RedisSecrets struct {
	Addr     string `json:"redis_addr"`
	Password string `json:"redis_port"`
	DB       int    `json:"redis_db"`
}

func (r RedisSecrets) Validate() error {
	var errs errsx.Map
	if r.Addr == "" {
		errs.Set("redis address", "redis address cannot be empty")
	}
	if r.Password == "" {
		errs.Set("redis password", "redis password must be between 1 and 65535")
	}
	if r.DB < 0 {
		errs.Set("redis DB", "redis DB must be non-negative")
	}
	return errs.AsError()
}

func (r RedisSecrets) GetType() string {
	return "redis"
}

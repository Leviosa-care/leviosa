package main

import (
	"fmt"
	"os"
	"strconv"
)

func setupEnvVars() error {
	// setup port or use default value 3500
	port_env := os.Getenv("BACKEND_PORT")
	if port_env == "" {
		port_env = "3500"
	}
	port_int, err := strconv.Atoi(port_env)
	if err != nil {
		return fmt.Errorf("port variable: %w: expect int", err)
	}
	opts.server.port = port_int

	// setup mode or use default value "development"
	if err = opts.mode.Set(os.Getenv("APP_ENV")); err != nil {
		return fmt.Errorf("mode: %w", err)
	}
	return nil
}

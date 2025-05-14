package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hengadev/leviosa/pkg/flags"

	"github.com/hengadev/errsx"
)

type sqliteCreds struct {
	Filename string `json:"filename"`
}

func (c *Config) GetSQLITE() *sqliteCreds {
	return c.sqlite
}

func (c *Config) setSQLITE(env mode.EnvMode) error {
	var errs errsx.Map
	databaseFilename := c.viper.GetString("sqlite.filename")
	if databaseFilename == "" {
		errs.Set("DATABASE_FILENAME", "'DATABASE_FILENAME' environment variable not set; please define it to specify SQLite file name")
	}
	var prefix string
	switch env {
	case mode.ModeStaging, mode.ModeDev, mode.ModeProd:
		prefix = env.String()
	default:
		errs.Set("mode value", fmt.Errorf("mode value can only be 'development', 'production' or 'staging', got : %q"))
	}
	sqliteFile := fmt.Sprintf("%s_%s", prefix, databaseFilename)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory to set SQLite database: %w", err)
	}
	c.sqlite.Filename = filepath.Join(wd, "data", sqliteFile)
	return nil
}

package logger

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/hengadev/leviosa/pkg/envmode"
)

func SetHandler(level, style string) (slog.Handler, error) {
	logLevel, ok := loggerLevels[loggerLevel(level)]
	if !ok {
		return nil, fmt.Errorf("invalid log level supplied: %q", level)
	}
	logStyle := loggerStyle(style)
	var slogHandler slog.Handler
	switch logStyle {
	case JSON:
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	case Text:
		slogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	case Dev:
		slogHandler = NewDevHandler(os.Stdout, logLevel)
	default:
		return nil, fmt.Errorf("invalid log style supplied: %q", style)
	}
	return slogHandler, nil
}

func SetOptions(env envmode.Mode, level, style *string) error {
	if env == envmode.Dev {
		*level = string(Debug)
		*style = string(Dev)
		return nil
	}
	var defaultLevel string
	switch env {
	case envmode.Prod:
		defaultLevel = string(Info)
	case envmode.Staging:
		defaultLevel = string(Debug)
	case envmode.Dev:
		defaultLevel = string(Debug)
	default:
		return fmt.Errorf("APP_ENV does not exist")
	}

	flag.StringVar(level, "logger-level", defaultLevel, "Set logger level")
	flag.StringVar(style, "logger-style", string(JSON), "Set logger style")
	flag.Parse()

	return nil
}

package logger

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"
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
	// Set defaults based on environment
	var defaultLevel, defaultStyle string
	switch env {
	case envmode.Dev:
		defaultLevel = string(Debug)
		defaultStyle = string(Dev)
	case envmode.Staging:
		defaultLevel = string(Debug)
		defaultStyle = string(JSON)
	case envmode.Prod:
		defaultLevel = string(Info)
		defaultStyle = string(JSON)
	default:
		return fmt.Errorf("invalid environment: %v", env)
	}

	// Only set up flags if not already parsed (avoid conflicts)
	if !flag.Parsed() {
		flag.StringVar(level, "logger-level", defaultLevel, "Set logger level (info, debug, error, warn)")
		flag.StringVar(style, "logger-style", defaultStyle, "Set logger style (json, text, dev)")
	} else {
		// Use defaults if flags already parsed
		if *level == "" {
			*level = defaultLevel
		}
		if *style == "" {
			*style = defaultStyle
		}
	}

	// Validate level and style combination
	return validateLoggerConfig(*level, *style)
}

// validateLoggerConfig validates logger configuration
func validateLoggerConfig(level, style string) error {
	// Validate level
	if _, ok := loggerLevels[loggerLevel(level)]; !ok {
		return fmt.Errorf("invalid log level: %q (valid: info, debug, error, warn)", level)
	}

	// Validate style
	validStyles := []string{string(JSON), string(Text), string(Dev)}
	if !slices.Contains(validStyles, style) {
		return fmt.Errorf("invalid log style: %q (valid: json, text, dev)", style)
	}

	return nil
}

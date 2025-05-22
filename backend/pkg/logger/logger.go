package logger

import (
	"log/slog"
)

type loggerLevel string
type loggerStyle string

const (
	Info  loggerLevel = "info"
	Debug             = "debug"
	Error             = "error"
	Warn              = "warn"
)

const (
	JSON loggerStyle = "json"
	Text             = "text"
	Dev              = "dev"
)

var loggerLevels = map[loggerLevel]slog.Level{
	Info:  slog.LevelInfo,
	Debug: slog.LevelDebug,
	Error: slog.LevelError,
	Warn:  slog.LevelWarn,
}

package logger

import (
	"log/slog"
)

type loggerLevel string
type loggerStyle string

const (
	Info  loggerLevel = "info"
	Debug loggerLevel = "debug"
	Error loggerLevel = "error"
	Warn  loggerLevel = "warn"
)
const (
	JSON loggerStyle = "json"
	Text loggerStyle = "text"
	Dev  loggerStyle = "dev"
)

var loggerLevels = map[loggerLevel]slog.Level{
	Info:  slog.LevelInfo,
	Debug: slog.LevelDebug,
	Error: slog.LevelError,
	Warn:  slog.LevelWarn,
}

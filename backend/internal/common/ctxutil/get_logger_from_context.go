package ctxutil

import (
	"context"
	"fmt"
	"log/slog"
)

// GetLoggerFromContext retrieves the custom logger from context using pseudonomyzation
func GetLoggerFromContext(ctx context.Context) (*slog.Logger, error) {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve logger from context: key %v is missing or value is not of type *slog.Logger", LoggerKey)
	}
	return logger, nil
}

package ctxutil

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

type contextKey string

const RoleKey = contextKey("role")
const UserIDKey = contextKey("userID")
const LoggerKey = contextKey("logger")

// GetLoggerFromContext retrieves the custom logger from context using pseudonomyzation
func GetLoggerFromContext(ctx context.Context) (*slog.Logger, error) {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve logger from context: key %v is missing or value is not of type *slog.Logger", LoggerKey)
	}
	return logger, nil
}

func ValidateRoleInContext(ctx context.Context, expectedRole models.Role) error {
	role, ok := ctx.Value(RoleKey).(models.Role)
	if !ok {
		return fmt.Errorf("role not found in context")
	}
	if role != expectedRole {
		return fmt.Errorf("expected role %q, got %q", expectedRole, role)
	}
	return nil
}

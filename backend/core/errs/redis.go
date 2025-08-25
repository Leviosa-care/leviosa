package errs

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// ClassifyRedisError maps specific Redis errors to your sentinel errors
func ClassifyRedisError(operation string, err error) error {
	if err == nil {
		return nil
	}

	var repoErr error
	switch {
	case errors.Is(err, redis.Nil):
		// Key not found
		repoErr = fmt.Errorf("%w: %w", ErrRepositoryNotFound, err)
	case errors.Is(err, redis.TxFailedErr):
		// Transaction failed
		repoErr = fmt.Errorf("%w: %w", ErrDatabase, err)
	case isRedisConnectionError(err):
		// Connection issues, timeouts, pool exhaustion
		repoErr = fmt.Errorf("%w: %w", ErrDBQuery, err)
	case isRedisContextError(err):
		// Context cancellation or timeout
		repoErr = fmt.Errorf("%w: %w", ErrContext, err)
	default:
		// For any other Redis error, wrap it with a general database error
		repoErr = fmt.Errorf("%w: %w", ErrDatabase, err)
	}

	// Always wrap the classified error with a clear message indicating the operation
	return fmt.Errorf("%s: %w", operation, repoErr)
}

// isRedisConnectionError checks if the error is related to connection issues
func isRedisConnectionError(err error) bool {
	errStr := err.Error()
	return errors.Is(err, redis.ErrClosed) ||
		containsAny(errStr, []string{
			"connection refused",
			"timeout",
			"pool exhausted",
			"broken pipe",
			"connection reset",
			"no route to host",
		})
}

// isRedisContextError checks if the error is context-related
func isRedisContextError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled)
}

// containsAny checks if the string contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				match := true
				for j := 0; j < len(sub); j++ {
					if s[i+j] != sub[j] {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}


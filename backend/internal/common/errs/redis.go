package errs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
)

// ClassifyRedisError maps specific Redis errors to your sentinel errors
func ClassifyRedisError(operation string, err error) error {
	if err == nil {
		return nil
	}

	// Handle context errors first
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", operation, ErrContext)
	}

	var repoErr error
	switch {
	case errors.Is(err, redis.Nil):
		// Key not found - this is normal for cache misses
		repoErr = fmt.Errorf("%w: %w", ErrRepositoryNotFound, err)

	case errors.Is(err, redis.TxFailedErr):
		// Redis transaction failed (MULTI/EXEC)
		repoErr = fmt.Errorf("%w: %w", ErrTransactionFailure, err)

	case errors.Is(err, redis.ErrClosed):
		// Redis client/connection is closed
		repoErr = fmt.Errorf("%w: %w", ErrConnectionFailure, err)

	// Network and connection errors
	case isNetworkError(err):
		// Network-level connection failures
		repoErr = fmt.Errorf("%w: %w", ErrConnectionFailure, err)

	case isRedisPoolError(err):
		// Connection pool exhaustion
		repoErr = fmt.Errorf("%w: %w", ErrTooManyConnections, err)

	case isRedisTimeoutError(err):
		// Redis operation timeout
		repoErr = fmt.Errorf("%w: %w", ErrQueryCancelled, err)

	case isRedisAuthError(err):
		// Redis authentication failure
		repoErr = fmt.Errorf("%w: %w", ErrPermissionDenied, err)

	case isRedisMemoryError(err):
		// Redis out of memory
		repoErr = fmt.Errorf("%w: %w", ErrResourceExhausted, err)

	default:
		// For any other Redis error, wrap it with a general database error
		repoErr = fmt.Errorf("%w: %w", ErrDatabase, err)
	}

	// Always wrap the classified error with a clear message indicating the operation
	return fmt.Errorf("%s: %w", operation, repoErr)
}

// isNetworkError checks if the error is a network-level connection failure
func isNetworkError(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return true
	}

	// Check for common network error strings as fallback
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "no route to host") ||
		strings.Contains(errStr, "network is unreachable")
}

// isRedisPoolError checks if the error is related to connection pool exhaustion
func isRedisPoolError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "pool exhausted") ||
		strings.Contains(errStr, "pool timeout") ||
		strings.Contains(errStr, "connection pool timeout")
}

// isRedisTimeoutError checks if the error is a Redis operation timeout
func isRedisTimeoutError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "read timeout") ||
		strings.Contains(errStr, "write timeout") ||
		strings.Contains(errStr, "dial timeout")
}

// isRedisAuthError checks if the error is Redis authentication related
func isRedisAuthError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "NOAUTH") ||
		strings.Contains(errStr, "invalid password") ||
		strings.Contains(errStr, "authentication failed") ||
		strings.Contains(errStr, "WRONGPASS")
}

// isRedisMemoryError checks if the error is Redis memory related
func isRedisMemoryError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "OOM") ||
		strings.Contains(errStr, "out of memory") ||
		strings.Contains(errStr, "maxmemory")
}

package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// RespondWithServiceError maps application/domain errors to appropriate HTTP responses.
// It handles error classification, logging, and response formatting in one place.
//
// Parameters:
//   - w: HTTP response writer
//   - logger: Structured logger for error logging
//   - ctx: Request context for contextual logging
//   - err: Error from application/service layer
//   - operation: Human-readable operation description (e.g., "create user", "update settings")
//
// Behavior:
//   - 4xx errors: Returns error message to client, no logging (client error)
//   - 5xx errors: Logs error with full context, returns generic message to client (server error)
func RespondWithServiceError(w http.ResponseWriter, logger *slog.Logger, ctx context.Context, err error, operation string) {
	if err == nil {
		return
	}

	// Classify error to HTTP status code
	statusCode := classifyServiceError(err)

	// Log server errors (5xx) - these indicate problems on our side
	if statusCode >= 500 {
		logger.ErrorContext(ctx,
			fmt.Sprintf("Handler: %s failed", operation),
			"error", err,
			"status_code", statusCode,
		)

		// Return generic message to client (don't leak internal details)
		RespondWithError(w, errors.New("an internal server error occurred"), statusCode)
		return
	}

	// Client errors (4xx) - return error message as-is
	RespondWithError(w, err, statusCode)
}

// classifyServiceError maps sentinel errors to HTTP status codes following REST best practices.
//
// Status Code Categories:
//   - 400-499: Client errors (invalid input, not found, unauthorized, etc.)
//   - 500-599: Server errors (database failures, internal errors, etc.)
//   - 503: Service Unavailable (retryable errors - connection failures, deadlocks, etc.)
//   - 408/504: Timeout errors (non-retryable timeout conditions)
func classifyServiceError(err error) int {
	switch {
	// ============================================================================
	// CLIENT ERRORS (4xx) - Problems with the request
	// ============================================================================

	// 400 Bad Request - Invalid input data or business validation failure
	case errors.Is(err, errs.ErrInvalidValue),
		errors.Is(err, errs.ErrInvalidInput),
		errors.Is(err, errs.ErrValidation),
		errors.Is(err, errs.ErrNoFieldsForUpdate):
		return http.StatusBadRequest

	// 401 Unauthorized - Authentication required or failed
	case errors.Is(err, errs.ErrUnauthorized),
		errors.Is(err, errs.ErrExpiredToken):
		return http.StatusUnauthorized

	// 403 Forbidden - Authenticated but not authorized for this resource
	case errors.Is(err, errs.ErrForbidden),
		errors.Is(err, errs.ErrPermissionDenied):
		return http.StatusForbidden

	// 404 Not Found - Resource does not exist
	case errors.Is(err, errs.ErrDomainNotFound),
		errors.Is(err, errs.ErrRepositoryNotFound):
		return http.StatusNotFound

	// 408 Request Timeout - Client took too long or cancelled request
	case errors.Is(err, context.Canceled),
		errors.Is(err, errs.ErrQueryCancelled):
		return http.StatusRequestTimeout

	// 409 Conflict - Resource already exists or conflict with current state
	case errors.Is(err, errs.ErrAlreadyExists),
		errors.Is(err, errs.ErrConflict),
		errors.Is(err, errs.ErrUniqueViolation):
		return http.StatusConflict

	// 423 Locked - Resource is locked (e.g., account locked)
	case errors.Is(err, errs.ErrAccountLocked):
		return http.StatusLocked

	// 429 Too Many Requests - Rate limit exceeded
	case errors.Is(err, errs.ErrRateLimit):
		return http.StatusTooManyRequests

	// ============================================================================
	// SERVER ERRORS (5xx) - Problems on our side
	// ============================================================================

	// 500 Internal Server Error - Generic server error
	case errors.Is(err, errs.ErrInternal),
		errors.Is(err, errs.ErrUnexpectedError),
		errors.Is(err, errs.ErrDatabase),
		errors.Is(err, errs.ErrDBQuery),
		errors.Is(err, errs.ErrQueryFailed),
		errors.Is(err, errs.ErrNotEncrypted),
		errors.Is(err, errs.ErrNotDecrypted),
		errors.Is(err, errs.ErrDomainNotCreated),
		errors.Is(err, errs.ErrDomainNotUpdated),
		errors.Is(err, errs.ErrDomainNotDeleted),
		errors.Is(err, errs.ErrRepositoryNotCreated),
		errors.Is(err, errs.ErrRepositoryNotUpdated),
		errors.Is(err, errs.ErrRepositoryNotDeleted):
		return http.StatusInternalServerError

	// 502 Bad Gateway - External service returned invalid response
	case errors.Is(err, errs.ErrExternalService):
		return http.StatusBadGateway

	// 503 Service Unavailable - Temporary failure, client should retry
	// These are RETRYABLE errors (connection pools, deadlocks, resource exhaustion)
	case errors.Is(err, errs.ErrConnectionFailure),
		errors.Is(err, errs.ErrTooManyConnections),
		errors.Is(err, errs.ErrResourceExhausted),
		errors.Is(err, errs.ErrTransactionFailure),
		errors.Is(err, errs.ErrDeadlock):
		return http.StatusServiceUnavailable

	// 504 Gateway Timeout - Upstream service timed out
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout

	// ============================================================================
	// DEFAULT - Unclassified errors
	// ============================================================================
	default:
		// Log warning about unhandled error type (should not happen in production)
		// This helps identify missing error classifications
		return http.StatusInternalServerError
	}
}

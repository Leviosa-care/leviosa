package aggregatorHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CheckEmailSendOTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload domain.CheckEmailAvailabilityRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "decode_request",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Log incoming request with masked email and request context
	logger.InfoContext(ctx, "Handler: Processing email verification request",
		"email", maskEmail(payload.Email),
		"operation", "check_email_send_otp",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	if err := h.svc.CheckEmailSendOTP(ctx, &payload); err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid email validation"
		case errors.Is(err, errs.ErrConflict):
			logLevel = "info"
			errorContext = "email already registered"
		case errors.Is(err, errs.ErrRateLimit):
			logLevel = "warn"
			errorContext = "OTP rate limit exceeded"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "infrastructure connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "infrastructure resource exhaustion"
		case errors.Is(err, errs.ErrExternalService):
			logLevel = "error"
			errorContext = "external service failure"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure):
			logLevel = "error"
			errorContext = "infrastructure transaction failure"
		default:
			logLevel = "error"
			errorContext = "unexpected error"
		}

		logFields := []any{
			"email", maskEmail(payload.Email),
			"operation", "check_email_send_otp",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", getStatusCodeForError(err),
			"error", err,
		}

		switch logLevel {
		case "info":
			logger.InfoContext(ctx, "Handler: Email verification request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Email verification request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Email verification request failed", logFields...)
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// Invalid input data
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrConflict):
			// Email already registered
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrRateLimit):
			// Rate limit exceeded
			statusCode = http.StatusTooManyRequests
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			// Infrastructure connection issues
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			// Infrastructure resources exhausted
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrExternalService):
			// External service (RabbitMQ) failure
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			// Query or request cancelled
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			// Request timed out
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure):
			// Temporary infrastructure issues - client should retry
			statusCode = http.StatusServiceUnavailable
		default:
			// Any other error
			statusCode = http.StatusInternalServerError
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation with enhanced context
	logger.InfoContext(ctx, "Handler: Email verification request completed successfully",
		"email", maskEmail(payload.Email),
		"operation", "check_email_send_otp",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Verification email sent successfully",
		Status:  "sent",
	}, http.StatusOK)
}

// getStatusCodeForError returns HTTP status code for given error (for logging)
func getStatusCodeForError(err error) int {
	switch {
	case errors.Is(err, errs.ErrInvalidValue):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, errs.ErrRateLimit):
		return http.StatusTooManyRequests
	case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
		return http.StatusServiceUnavailable
	case errors.Is(err, errs.ErrResourceExhausted):
		return http.StatusServiceUnavailable
	case errors.Is(err, errs.ErrExternalService):
		return http.StatusServiceUnavailable
	case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
		return http.StatusRequestTimeout
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusRequestTimeout
	case errors.Is(err, errs.ErrTransactionFailure):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

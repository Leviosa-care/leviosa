package aggregatorHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware/auth"
)

func (h *handler) ValidateOTPCreatePendingUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload domain.ValidateOTPRequest

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

	logger.InfoContext(ctx, "Handler: Processing OTP verification request",
		"email", maskEmail(payload.Email),
		"operation", "validate_otp_create_pending_user",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	session, err := h.svc.ValidateOTPCreatePendingUser(ctx, &payload)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "validation failed"
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "resource not found"
		case errors.Is(err, errs.ErrExpiredToken):
			logLevel = "info"
			errorContext = "token expired"
		case errors.Is(err, errs.ErrRateLimit):
			logLevel = "warn"
			errorContext = "rate limit exceeded"
		case errors.Is(err, errs.ErrValueMismatch):
			logLevel = "warn"
			errorContext = "authentication failed"
		case errors.Is(err, errs.ErrAlreadyConsumed):
			logLevel = "warn"
			errorContext = "resource already consumed"
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
			"operation", "validate_otp_create_pending_user",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// Invalid OTP format or email validation
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrDomainNotFound):
			// OTP doesn't exist or already used
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrExpiredToken):
			// OTP has expired
			statusCode = http.StatusGone
		case errors.Is(err, errs.ErrRateLimit):
			// Maximum attempts exceeded
			statusCode = http.StatusTooManyRequests
		case errors.Is(err, errs.ErrValueMismatch):
			// Wrong OTP code
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrAlreadyConsumed):
			// OTP already consumed by concurrent request
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			// Infrastructure connection issues
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			// Infrastructure resources exhausted
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrExternalService):
			// External service failure
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

		logFields = append(logFields, "status_code", statusCode)

		switch logLevel {
		case "info":
			logger.InfoContext(ctx, "Handler: OTP validation request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: OTP validation request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: OTP validation request failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation with enhanced context
	logger.InfoContext(ctx, "Handler: OTP verification request completed successfully",
		"email", maskEmail(payload.Email),
		"operation", "validate_otp_create_pending_user",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusCreated)

	// Set dual token cookies
	mw.SetTokenCookies(w, session.AccessToken, session.RefreshToken,
		session.AccessTokenExpiry, session.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Pending user created successfully",
		Status:  "created",
	}, http.StatusCreated)
}

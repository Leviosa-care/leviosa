package aggregatorHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CompletePartner(w http.ResponseWriter, r *http.Request) {
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

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.WarnContext(ctx, "Handler: Missing session info",
			"error", err,
			"operation", "complete_partner",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
		return
	}

	// Parse request body
	var payload domain.CompletePartnerRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "complete_partner",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing complete partner request",
		"operation", "complete_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call the aggregator service
	err = h.svc.CompletePartner(ctx, sessionInfo, &payload)
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
			errorContext = "session not found"
		case errors.Is(err, errs.ErrExpiredToken):
			logLevel = "info"
			errorContext = "session expired"
		case errors.Is(err, errs.ErrConflict):
			logLevel = "warn"
			errorContext = "invalid state or conflict"
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
			"operation", "complete_partner",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}

		var statusCode int
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// Invalid request data or validation failed
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrDomainNotFound):
			// Session doesn't exist
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrExpiredToken):
			// Session has expired
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errs.ErrConflict):
			// Session/user state conflict or already completed
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			// Infrastructure connection issues
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			// Infrastructure resources exhausted
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrExternalService):
			// External service failure (e.g., Stripe API, catalog service)
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
			logger.InfoContext(ctx, "Handler: Complete partner request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Complete partner request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Complete partner request failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Complete partner request completed successfully",
		"operation", "complete_partner",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "Partner registration completed successfully. Awaiting admin approval.",
	}, http.StatusOK)
}

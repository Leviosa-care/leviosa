package aggregatorHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

func (h *handler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get session info from context (injected by auth middleware)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: No session info in context",
			"operation", "sign_out",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing sign-out request",
		"operation", "sign_out",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call the aggregator service to sign out the user
	err = h.svc.SignOut(ctx, sessionInfo)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
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
			"operation", "sign_out",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Sign-out request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Sign-out request failed", logFields...)
		}

		var statusCode int
		switch {
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

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Sign-out request completed successfully",
		"operation", "sign_out",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Successfully signed out",
		Status:  "signed_out",
	}, http.StatusOK)
}


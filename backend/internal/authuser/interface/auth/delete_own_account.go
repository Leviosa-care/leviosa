package aggregatorHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) DeleteOwnAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract logger from context
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get session info from context (user ID comes from here)
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("session information from context required"), http.StatusUnauthorized)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete own account request",
		"operation", "delete_own_account",
		"user_id", sessionInfo.UserID.String(),
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Call the aggregator service to delete the user's own account
	err = h.svc.DeleteOwnAccount(ctx, sessionInfo)
	if err != nil {
		var statusCode int
		var errorContext string
		var logLevel string

		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
			errorContext = "account not found"
			logLevel = "warn"
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
			errorContext = "validation failed"
			logLevel = "warn"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database connection issue"
			logLevel = "error"
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
			errorContext = "resource exhausted"
			logLevel = "error"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
			logLevel = "warn"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
			logLevel = "warn"
		case errors.Is(err, errs.ErrTransactionFailure):
			statusCode = http.StatusServiceUnavailable
			errorContext = "transaction failure"
			logLevel = "error"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "internal server error"
			logLevel = "error"
		}

		logFields := []any{
			"operation", "delete_own_account",
			"user_id", sessionInfo.UserID.String(),
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", statusCode,
			"error", err,
		}

		switch logLevel {
		case "warn":
			logger.WarnContext(ctx, "Handler: Delete own account failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Delete own account failed", logFields...)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful completion
	logger.InfoContext(ctx, "Handler: Delete own account completed",
		"operation", "delete_own_account",
		"user_id", sessionInfo.UserID.String(),
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "Account deleted successfully",
	}, http.StatusOK)
}
